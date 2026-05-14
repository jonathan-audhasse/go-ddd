package transaction_test

import (
	"context"
	"goddd/internal/infrastructure/persistence/testutil"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	testutil.InitTestLogger()

	ctx := context.Background()

	pgCon, err := testutil.NewPGContainer(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed create new pg container")
	}

	// close pg container
	defer func() {
		pgCon.Close(ctx)
	}()

	// set DB
	testDB = pgCon.Pool

	code := m.Run()

	os.Exit(code)
}
