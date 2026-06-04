package postgres

import (
	"goddd/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewRepositories instantiates a new Repositories
func NewRepositories(pool *pgxpool.Pool) *repository.Repositories {
	return &repository.Repositories{
		UserRepo: NewUserRepository(pool),
	}
}
