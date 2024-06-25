package postgres

import (
	"goddd/src/domain/repository"

	"github.com/jmoiron/sqlx"
)

// Factory
func NewRepository(db *sqlx.DB) *repository.Repository {
	return &repository.Repository{CustomerRepo: NewCustomerPgRepo(db)}
}
