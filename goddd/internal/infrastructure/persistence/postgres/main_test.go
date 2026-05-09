package postgres_test

import (
	"context"
	"fmt"
	"goddd/internal/infrastructure/persistence/migrations"
	"goddd/internal/infrastructure/persistence/postgres"
	"goddd/internal/infrastructure/persistence/transaction"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	migrationsPath = "file://../migrations/files"
	testSchemaName = "test_schema"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	// init logger for test
	// to see logs add the following line in the test you want to debug
	// ctx := log.Logger.WithContext(context.Background())
	// or
	// ctx := log.Logger.WithContext(t.Context())
	initTestLogger()

	ctx := context.Background()

	pgCon, err := newPGContainer(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed create new pg container")
	}

	// close pg container
	defer func() {
		pgCon.Close(ctx)
	}()

	// set DB
	testDB = pgCon.pool

	code := m.Run()

	os.Exit(code)
}

// pgContainer is a wrapper around the test container
type pgContainer struct {
	container *tcpostgres.PostgresContainer
	pool      *pgxpool.Pool
	dbUrl     string
}

// newPGContainer create a new postgres container
func newPGContainer(ctx context.Context) (*pgContainer, error) {
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

	pool, err := postgres.Open(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test DB: %w", err)
	}

	// run migration
	migrations.Up(pool, migrationsPath, testSchemaName)

	return &pgContainer{
		container: container,
		pool:      pool,
		dbUrl:     dsn,
	}, nil
}

// Terminate stops the container
func (c pgContainer) Close(ctx context.Context) {
	c.pool.Close()
	if err := c.container.Terminate(ctx); err != nil {
		log.Error().Err(err).Msg("failed to terminate postgres container")
	}
}

// initTestLogger for debugging
func initTestLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout},
	).With().Timestamp().Logger()

	log.Logger = logger
}

// withTx helps having isolation per test.
// rollback transactions per test
func withTx(t *testing.T, ctx context.Context, db *pgxpool.Pool) (context.Context, pgx.Tx, func()) {
	t.Helper()

	tx, err := db.Begin(ctx)
	require.NoError(t, err)

	// attach tx to context
	ctx = transaction.WithTx(ctx, tx)

	rollback := func() {
		_ = tx.Rollback(ctx)
	}

	return ctx, tx, rollback
}
