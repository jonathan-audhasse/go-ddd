package bootstrap

import (
	"goddd/internal/application"
	"goddd/internal/domain/repository"
	"goddd/internal/infrastructure/persistence/postgres"
	txpg "goddd/internal/infrastructure/persistence/transaction"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LoadServices loads services
func LoadServices(db *pgxpool.Pool, repos *repository.Repositories) *application.Services {

	tm := txpg.NewTransactionManager(db)

	pinger := postgres.NewPinger(db)

	return application.NewServices(repos, tm, pinger)
}
