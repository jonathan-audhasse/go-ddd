package healthservice

import (
	"context"
	"goddd/internal/domain/repository"
)

// HealthService for asset's services implementation
type HealthService struct {
	repo repository.HealthRepository
}

// Create the HealthService
func NewHealthService(repo repository.HealthRepository) *HealthService {
	return &HealthService{repo}
}

// IsHealthy health check
func (hs *HealthService) IsHealthy(ctx context.Context) error {
	// ping repo
	if err := hs.repo.Ping(ctx); err != nil {
		return err
	}
	// healthy
	return nil
}
