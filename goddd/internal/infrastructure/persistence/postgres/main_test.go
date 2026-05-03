package postgres_test

import (
	"context"
	"goddd/internal/infrastructure/persistence/postgres"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	ctx := context.Background()

	container, err := tcpostgres.Run(
		ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start postgres container")
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get postgres connection string")
	}

	db, err := postgres.Open(ctx, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to test DB")
	}

	testDB = db

	code := m.Run()

	db.Close()
	if err := container.Terminate(ctx); err != nil {
		log.Error().Err(err).Msg("failed to terminate postgres container")
	}

	os.Exit(code)
}

// initTestLogger for debugging
func initTestLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout},
	).With().Timestamp().Logger()

	log.Logger = logger
}

// withTx helps having isolation per test.
// rollback transactions per test
func withTx(t *testing.T, db *pgxpool.Pool) (pgx.Tx, func()) {
	t.Helper()

	tx, err := db.Begin(context.Background())
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}

	rollback := func() {
		_ = tx.Rollback(context.Background())
	}

	return tx, rollback
}
