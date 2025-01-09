package postgres

import (
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type customerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) *customerRepo {
	return &customerRepo{db}
}

// Get a customer by ID
func (r *customerRepo) Get(id string) (models.Customer, error) {
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
func (r *customerRepo) Add(cust models.Customer) error {
	tx := r.db.MustBegin()
	if _, err := tx.NamedExec("INSERT INTO customer (id, user_id, name, email) VALUES (:id, :user_id, :name, :email)", &cust); err != nil {
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
func (r *customerRepo) ListByUser(userid string) ([]models.Customer, error) {
	if userid == "" {
		return []models.Customer{}, errors.InternalError.New("user must be specified")
	}
	cust := make([]models.Customer, 0)
	if err := r.db.Select(&cust, "SELECT * FROM customer WHERE user_id=$1", userid); err != nil {
		log.Printf("failed to execute queyr for user id=%s: %s\n", userid, err)
		return cust, errors.InternalError.Wrap(err, "fail to list customer from repository")
	}
	return cust, nil
}

// Delete a customer by ID
func (r *customerRepo) Delete(id string) error {
	if id == "" {
		return errors.InternalError.New("customer id must not be empty")
	}
	if _, err := r.db.Exec("DELETE FROM customer WHERE id=$1", id); err != nil {
		return errors.InternalError.Wrapf(err, "fail to delete customer id=%v from repository", id)
	}
	return nil
}
