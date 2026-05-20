package bootstrap

import (
	"goddd/internal/application"
	"goddd/internal/domain/repository"
	txpg "goddd/internal/infrastructure/persistence/transaction"

	"github.com/jackc/pgx/v5/pgxpool"
)

func LoadServices(db *pgxpool.Pool, repos *repository.Repositories) *application.Services {

	tm := txpg.NewTransactionManager(db)

	return application.NewServices(repos, tm)
}
