package postgres

import (
	"goddd/src/domain/models"

	"github.com/jmoiron/sqlx"
)

type CustomerPgRepo struct {
	db *sqlx.DB
}

func NewCustomerPgRepo(db *sqlx.DB) *CustomerPgRepo {
	return &CustomerPgRepo{db}
}

// Get a customer by ID
func (r *CustomerPgRepo) Get(id string) (models.Customer, error) {
	return models.Customer{}, nil
}

// Add a new customer to the repository
func (r *CustomerPgRepo) Add(cust models.Customer) error {
	return nil
}

// List customers
func (r *CustomerPgRepo) List() ([]models.Customer, error) {
	return []models.Customer{}, nil
}
