package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type HealthRepository interface {
	CheckHealth(ctx context.Context) (*HealthStatus, error)
}

type healthRepository struct {
	db      *sql.DB
	version string
}

func NewHealthRepository(db *sql.DB, version string) HealthRepository {
	return &healthRepository{
		db:      db,
		version: version,
	}
}

func (r *healthRepository) CheckHealth(ctx context.Context) (*HealthStatus, error) {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Version:   r.version,
	}

	if err := r.db.PingContext(ctx); err != nil {
		status.Status = "unhealthy"
		status.DB = "down"
		return status, fmt.Errorf("database connection failed: %w", err)
	}

	status.DB = "up"
	return status, nil
}
