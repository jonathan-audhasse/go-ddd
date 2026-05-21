package e2e_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"goddd/internal/bootstrap"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testServer *http.Server
	testDB     *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	app, cleanup, err := bootstrap.NewTestApplication(ctx)
	if err != nil {
		panic(err)
	}

	testServer = app.Server
	testDB = app.DB

	code := m.Run()

	cleanup()

	os.Exit(code)
}
