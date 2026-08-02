package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	ErrOpenConnection = errors.New("failed to open a connection")
	ErrPingDatabase   = errors.New("failed to ping database")
)

// Open creates and validates a pgxpool connection pool.
// schema controls PostgreSQL search_path.
// Example:
//   - public (dev)
//   - e2e (tests)
func Open(
	ctx context.Context,
	dsn string,
	schema string,
) (*pgxpool.Pool, error) {

	logger := log.With().
		Str("dataSourceName", dsn).
		Str("schema", schema).
		Logger()

	logger.Info().Msg("opening connection to database")

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Err(err).Msg("failed to parse database config")
		return nil, ErrOpenConnection
	}

	// configure PostgreSQL schema
	if schema != "" {
		cfg.ConnConfig.RuntimeParams["search_path"] = schema
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Err(err).Msg("failed to open connection")
		return nil, ErrOpenConnection
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Err(err).Msg("failed to ping database")
		return nil, ErrPingDatabase
	}

	logger.Info().Msg("database connection established")

	return pool, nil
}
