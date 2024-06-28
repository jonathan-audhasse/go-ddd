package memory

import (
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"sync"
)

// customerMemRepo fulfills the CustomerRepository interface
type customerMemRepo struct {
	customers map[string]models.Customer
	sync.Mutex
}

// New is a factory function to generate a new repository of customers
func NewCustomerMemRepo() *customerMemRepo {
	return &customerMemRepo{customers: make(map[string]models.Customer)}
}

// Get a customer by ID
func (r *customerMemRepo) Get(id string) (models.Customer, error) {
	if cust, ok := r.customers[id]; ok {
		return cust, nil
	}
	return models.Customer{}, errors.RepoItemNotFound.Newf("customer id=%v not found in repository", id)
}

// Add a new customer to the repository
func (r *customerMemRepo) Add(cust models.Customer) error {
	if r.customers == nil {
		r.Lock()
		r = &customerMemRepo{customers: make(map[string]models.Customer)}
		r.Unlock()
	}
	if _, ok := r.customers[cust.Id]; ok {
		return errors.InternalError.Newf("fail to add a new customer %v", cust)
	}
	r.Lock()
	r.customers[cust.Id] = cust
	r.Unlock()
	return nil
}

// List customers
func (r *customerMemRepo) List() ([]models.Customer, error) {
	cc := make([]models.Customer, 0)
	for _, cust := range r.customers {
		cc = append(cc, cust)
	}
	return cc, nil
}

// Remove a customer from repository
func (r *customerMemRepo) Delete(id string) error {
	delete(r.customers, id)
	return nil
}
