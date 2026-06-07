package health

import (
	"context"
	"fmt"
	"goddd/internal/application/apperror"
)

var ErrFailToCheckHealth = fmt.Errorf("%w: failed to check health", apperror.ErrUnavailable)

// HealthService Interface
type HealthService interface {
	IsHealthy(context.Context) error
}

// Pinger a service to check repository connection
type Pinger interface {
	Ping(context.Context) error
}

// service for asset's services implementation
type service struct {
	pinger Pinger
}

// Create the HealthService
func NewHealthService(pinger Pinger) HealthService {
	return &service{pinger}
}

// IsHealthy health check
func (hs *service) IsHealthy(ctx context.Context) error {

	// ping repo
	if err := hs.pinger.Ping(ctx); err != nil {
		return ErrFailToCheckHealth
	}
	// healthy
	return nil
}
