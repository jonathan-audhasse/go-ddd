package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// Open creates and validates a pgxpool connection pool.
// pgxpool is safe for concurrent use and is the recommended default for servers.
func Open(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	logger := log.With().Str("dsn", dsn).Logger()
	logger.Info().Msg("opening connection to database")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Err(err).Msg("failed to open connection")
		return nil, errors.New("failed to open a connection")
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, errors.New("failed to ping database")
	}
	return pool, nil
}
 