package postgres

import (
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type customerPgRepo struct {
	db *sqlx.DB
}

func NewCustomerPgRepo(db *sqlx.DB) *customerPgRepo {
	return &customerPgRepo{db}
}

// Get a customer by ID
func (r *customerPgRepo) Get(id string) (models.Customer, error) {
	if id == "" {
		return models.Customer{}, errors.InternalError.New("customer id must not be empty")
	}
	var cust models.Customer
	if err := r.db.Get(&cust, "SELECT * FROM customer WHERE id=$1", id); err != nil {
		return models.Customer{}, errors.RepoItemNotFound.Wrapf(err, "customer id=%v not found in repository", id)
	}
	return cust, nil
}

// Add a new customer to the repository
func (r *customerPgRepo) Add(cust models.Customer) error {
	tx := r.db.MustBegin()
	if _, err := tx.NamedExec("INSERT INTO customer (id, name, email) VALUES (:id, :name, :email)", &cust); err != nil {
		if e, ok := err.(*pq.Error); ok {
			switch e.Code.Name() {
			case "unique_violation":
				return errors.RepoItemAlreadyExist.Wrapf(e, "customer already exist in repository")
			default:
				return errors.InternalError.Wrapf(err, "fail to add customer (%v) from repository", cust)
			}
		}
		return errors.InternalError.Wrapf(err, "fail to add customer (%v) from repository", cust)
	}
	if err := tx.Commit(); err != nil {
		return errors.InternalError.Wrapf(err, "fail to add customer (%v) from repository", cust)
	}
	return nil
}

// List customers
func (r *customerPgRepo) List() ([]models.Customer, error) {
	cust := make([]models.Customer, 0)
	if err := r.db.Select(&cust, "SELECT * FROM customer"); err != nil {
		return cust, errors.InternalError.Wrap(err, "fail to list customer from repository")
	}
	return cust, nil
}

// Delete a customer by ID
func (r *customerPgRepo) Delete(id string) error {
	if id == "" {
		return errors.InternalError.New("customer id must not be empty")
	}
	if _, err := r.db.Exec("DELETE FROM customer where id=$1", id); err != nil {
		return errors.InternalError.Wrapf(err, "fail to delete customer id=%v from repository", id)
	}
	return nil
}
