package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-web-demo/internal/config"
	"go-web-demo/internal/handler"
	"go-web-demo/internal/repository"
	"go-web-demo/internal/router"
	"go-web-demo/internal/service"
	"go-web-demo/pkg/logger"
)

const Version = "1.0.0"

func main() {
	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log := logger.New(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.OutputPath,
	)

	log.Info("Starting application", "version", Version, "port", cfg.Server.Port)

	db, err := repository.InitDB(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize database", "error", err)
	}
	defer db.Close()

	app := setupApp(cfg, db, log)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      app.GetEngine(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	startServer(server, cfg.Server.ShutdownTimeout, log)
}

func setupApp(cfg *config.Config, db *sql.DB, log *logger.Logger) *router.Router {
	userRepo := repository.NewUserRepository(db)
	healthRepo := repository.NewHealthRepository(db, Version)

	userService := service.NewUserService(userRepo)
	healthService := service.NewHealthService(healthRepo)
	helloService := service.NewHelloService()

	handlers := &router.Handlers{
		Hello:  handler.NewHelloHandler(helloService),
		Health: handler.NewHealthHandler(healthService),
		User:   handler.NewUserHandler(userService),
	}

	routerCfg := &router.RouterConfig{
		Mode: cfg.Server.Mode,
	}

	app := router.New(routerCfg, log, handlers)
	app.Setup()

	return app
}

func startServer(server *http.Server, shutdownTimeout time.Duration, log *logger.Logger) {
	go func() {
		log.Info("Server is running", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server exited gracefully")
}
