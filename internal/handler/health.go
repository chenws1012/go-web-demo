package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-web-demo/internal/repository"
)

type HealthHandler struct {
	healthService HealthService
}

type HealthService interface {
	CheckHealth(ctx context.Context) (*repository.HealthStatus, error)
}

type HealthStatusResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	DB        string `json:"db"`
	Version   string `json:"version"`
}

func NewHealthHandler(healthService HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	status, err := h.healthService.CheckHealth(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	responseStatus := HealthStatusResponse{
		Status:    status.Status,
		Timestamp: status.Timestamp.Format(time.RFC3339),
		DB:        status.DB,
		Version:   status.Version,
	}

	c.JSON(http.StatusOK, responseStatus)
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
