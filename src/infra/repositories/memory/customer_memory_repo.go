package memory

import (
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"sync"
)

// CustomerMemRepo fulfills the CustomerRepository interface
type CustomerMemRepo struct {
	customers map[string]models.Customer
	sync.Mutex
}

// New is a factory function to generate a new repository of customers
func NewCustomerMemRepo() *CustomerMemRepo {
	return &CustomerMemRepo{customers: make(map[string]models.Customer)}
}

// Get a customer by ID
func (r *CustomerMemRepo) Get(id string) (models.Customer, error) {
	if cust, ok := r.customers[id]; ok {
		return cust, nil
	}
	return models.Customer{}, errors.RepoItemNotFound.Newf("customer (id=%v) not found in repository", id)
}

// Add a new customer to the repository
func (r *CustomerMemRepo) Add(cust models.Customer) error {
	if r.customers == nil {
		r.Lock()
		r = &CustomerMemRepo{customers: make(map[string]models.Customer)}
		r.Unlock()
	}
	if _, ok := r.customers[cust.ID]; ok {
		return errors.InternalError.Newf("fail to add a new customer %v", cust)
	}
	r.Lock()
	r.customers[cust.ID] = cust
	r.Unlock()
	return nil
}

// List customers
func (r *CustomerMemRepo) List() ([]models.Customer, error) {
	cc := make([]models.Customer, 0)
	for _, cust := range r.customers {
		cc = append(cc, cust)
	}
	return cc, nil
}

// Empty customers
func (r *CustomerMemRepo) Delete() error {
	for k := range r.customers {
		delete(r.customers, k)
	}
	return nil
}
