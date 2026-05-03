package postgres

import (
	"context"
	"errors"
	"goddd/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var ErrFailToPing = errors.New("failed to ping repository")

// healthRepository object to hold health repository operations
type healthRepository struct {
	db *pgxpool.Pool
}

// NewHealthRepository instanciate new health repository
func NewHealthRepository(pool *pgxpool.Pool) repository.HealthRepository {
	return &healthRepository{db: pool}
}

// Ping pings DB
func (r *healthRepository) Ping(ctx context.Context) error {
	log.Ctx(ctx).Debug().Msg("Jonathan")
	if err := r.db.Ping(ctx); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to ping postgres")
		return ErrFailToPing
	}

	return nil
}
