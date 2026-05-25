package testutils

import (
	"context"
	"fmt"
	"goddd/internal/infrastructure/persistence/migrations"
	"goddd/internal/infrastructure/persistence/postgres"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testSchemaName = "test_schema"
)

// PgContainer is a wrapper around the test container
type PgContainer struct {
	container *tcpostgres.PostgresContainer
	Pool      *pgxpool.Pool
}

// NewPGContainer create a new postgres container
func NewPGContainer(ctx context.Context) (*PgContainer, error) {
	container, err := tcpostgres.Run(
		ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpwd"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres connection string: %w", err)
	}

	pool, err := postgres.Open(ctx, dsn, "")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test DB: %w", err)
	}

	// run migration
	migrations.Up(dsn, testSchemaName)

	return &PgContainer{
		container: container,
		Pool:      pool,
	}, nil
}

// Terminate stops the container
func (c PgContainer) Close(ctx context.Context) {
	c.Pool.Close()
	if err := c.container.Terminate(ctx); err != nil {
		log.Error().Err(err).Msg("failed to terminate postgres container")
	}
}
