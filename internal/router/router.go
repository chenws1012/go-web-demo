package router

import (
	"github.com/gin-gonic/gin"
	"go-web-demo/internal/handler"
	"go-web-demo/internal/middleware"
	"go-web-demo/pkg/logger"
)

type Router struct {
	engine   *gin.Engine
	config   *RouterConfig
	log      *logger.Logger
	handlers *Handlers
}

type RouterConfig struct {
	Mode string
}

type Handlers struct {
	Hello  *handler.HelloHandler
	Health *handler.HealthHandler
	User   *handler.UserHandler
}

func New(cfg *RouterConfig, log *logger.Logger, handlers *Handlers) *Router {
	gin.SetMode(cfg.Mode)
	engine := gin.New()

	return &Router{
		engine:   engine,
		config:   cfg,
		log:      log,
		handlers: handlers,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(
		middleware.RequestID(),
		middleware.LoggerWithContext(r.log),
		middleware.Logging(r.log),
		middleware.Recovery(r.log),
		middleware.CORS(),
	)

	r.setupHealthRoutes()
	r.setupAPIRoutes()

	return r.engine
}

func (r *Router) setupHealthRoutes() {
	health := r.engine.Group("/health")
	{
		health.GET("/liveness", r.handlers.Health.Liveness)
		health.GET("/readiness", r.handlers.Health.CheckHealth)
	}
}

func (r *Router) setupAPIRoutes() {
	v1 := r.engine.Group("/api/v1")
	{
		v1.GET("/hello", r.handlers.Hello.SayHello)

		users := v1.Group("/users")
		{
			users.POST("", r.handlers.User.CreateUser)
			users.GET("", r.handlers.User.ListUsers)
			users.GET("/:id", r.handlers.User.GetUser)
			users.PUT("/:id", r.handlers.User.UpdateUser)
			users.DELETE("/:id", r.handlers.User.DeleteUser)
		}
	}
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
