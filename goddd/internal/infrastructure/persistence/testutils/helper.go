package testutils

import (
	"context"
	"goddd/internal/infrastructure/persistence/transaction"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

// initTestLogger for debugging
func InitTestLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout},
	).With().Timestamp().Logger()

	log.Logger = logger
}

// WithTx helps having isolation per test.
// rollback transactions per test
func WithTx(t *testing.T, ctx context.Context, db *pgxpool.Pool) (context.Context, pgx.Tx, func()) {
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
