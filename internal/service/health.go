package service

import (
	"context"

	"go-web-demo/internal/repository"
)

type HealthService interface {
	CheckHealth(ctx context.Context) (*repository.HealthStatus, error)
}

type healthService struct {
	healthRepo repository.HealthRepository
}

func NewHealthService(healthRepo repository.HealthRepository) HealthService {
	return &healthService{
		healthRepo: healthRepo,
	}
}

func (s *healthService) CheckHealth(ctx context.Context) (*repository.HealthStatus, error) {
	return s.healthRepo.CheckHealth(ctx)
}
