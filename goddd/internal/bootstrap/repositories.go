package bootstrap

import (
	"goddd/internal/domain/repository"
	"goddd/internal/infrastructure/persistence/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func LoadRepositories(db *pgxpool.Pool) *repository.Repositories {
	return postgres.NewRepositories(db)
}
