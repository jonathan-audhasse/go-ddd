package postgres

import (
	"context"
	"goddd/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// healthRepository object to hold health repository operations
type healthRepository struct {
	db *pgxpool.Pool
}

// NewHealthRepository instantiate new health repository
func NewHealthRepository(pool *pgxpool.Pool) repository.HealthRepository {
	return &healthRepository{db: pool}
}

// Ping pings DB
func (r *healthRepository) Ping(ctx context.Context) error {
	logger := log.Ctx(ctx)
	logger.Debug().Msg("try to ping database...")
	// log.Ctx(ctx).Debug().Msg("try to ping database...")
	if err := r.db.Ping(ctx); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to ping postgres")
		return repository.ErrFailToPing
	}

	logger.Debug().Msg("ping database succeed")

	return nil
}
