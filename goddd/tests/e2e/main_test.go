package e2e_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	httapi "goddd/api/http"
	"goddd/internal/application"
	"goddd/internal/bootstrap"
	"goddd/internal/infrastructure/config"
	"goddd/internal/infrastructure/logger"
	"goddd/internal/infrastructure/persistence/migrations"
	"goddd/internal/infrastructure/persistence/postgres"
	"goddd/internal/infrastructure/persistence/transaction"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const testSchemaName = "e2e_test"

var (
	testHandler http.Handler
	testDB      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// init logger
	logger.New(cfg.LogLevel)

	// open DB
	dsn := "postgres://admin:abc123@localhost:5432/postgres?sslmode=prefer"
	db, err := postgres.Open(ctx, dsn, testSchemaName)
	if err != nil {
		panic(err)
	}

	// migrate DB
	if err := migrations.Up(dsn, testSchemaName); err != nil {
		panic(err)
	}

	// build repositories
	repos := bootstrap.LoadRepositories(db)

	// transaction manager
	tm := transaction.NewTransactionManager(db)

	// health pinger
	pinger := postgres.NewPinger(db)

	// build services
	services := application.NewServices(repos, tm, pinger)

	// router
	r := httapi.NewRouter(services)

	testDB = db

	testHandler = r

	code := m.Run()

	// cleanup
	db.Close()

	os.Exit(code)
}

func resetDB(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	// order matters if you add more tables later
	_, err := testDB.Exec(ctx, `
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}
