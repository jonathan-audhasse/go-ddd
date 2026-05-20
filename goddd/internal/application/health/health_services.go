package health

import (
	"context"
	"goddd/internal/application/apperror"
	"goddd/internal/domain/repository"
)

// HealthService Interface
type HealthService interface {
	IsHealthy(context.Context) error
}

// service for asset's services implementation
type service struct {
	repo repository.HealthRepository
}

// Create the HealthService
func NewHealthService(repo repository.HealthRepository) HealthService {
	return &service{repo}
}

// IsHealthy health check
func (hs *service) IsHealthy(ctx context.Context) error {

	// ping repo
	if err := hs.repo.Ping(ctx); err != nil {
		return apperror.ToAppError(err)
	}
	// healthy
	return nil
}
