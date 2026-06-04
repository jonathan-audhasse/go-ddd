package postgres

import (
	"context"
	"errors"
	"goddd/internal/application/health"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var ErrFailToPing = errors.New("failed to ping repository")

// healthPinger object to hold health repository operations
type healthPinger struct {
	db *pgxpool.Pool
}

// NewPinger instantiate new health pinger
func NewPinger(pool *pgxpool.Pool) health.Pinger {
	return &healthPinger{db: pool}
}

// Ping pings DB
func (r *healthPinger) Ping(ctx context.Context) error {
	logger := log.Ctx(ctx)
	logger.Debug().Msg("try to ping database...")
	if err := r.db.Ping(ctx); err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to ping postgres")
		return ErrFailToPing
	}

	logger.Debug().Msg("ping database succeed")

	return nil
}
