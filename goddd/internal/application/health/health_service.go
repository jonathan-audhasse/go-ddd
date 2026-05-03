package health

import (
	"context"
	"goddd/internal/domain/repository"

	"github.com/rs/zerolog/log"
)

// HealthService for asset's services implementation
type service struct {
	repo repository.HealthRepository
}

// NewService create the health service
func NewService(repo repository.HealthRepository) *service {
	return &service{repo}
}

// IsReady check if the database is mounted.
// return true if database is reachable, false otherwise
func (s *service) IsReady(ctx context.Context) bool {

	if err := s.repo.Ping(ctx); err != nil {
		log.Ctx(ctx).Err(err).Msg("database is not ready")
		return false
	}

	log.Ctx(ctx).Info().Msg("database mounted and ready")
	return true
}
