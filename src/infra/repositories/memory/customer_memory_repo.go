package memory

import (
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"sync"
)

// CustomerMemoryRepository fulfills the CustomerRepository interface
type CustomerMemoryRepository struct {
	customers map[string]models.Customer
	sync.Mutex
}

// New is a factory function to generate a new repository of customers
func NewCustomerMemoryRepository() *CustomerMemoryRepository {
	return &CustomerMemoryRepository{customers: make(map[string]models.Customer)}
}

// Get a customer by ID
func (mr *CustomerMemoryRepository) Get(id string) (models.Customer, error) {
	if cust, ok := mr.customers[id]; ok {
		return cust, nil
	}
	return models.Customer{}, errors.RepoItemNotFound.Newf("customer (id=%v) not found in repository", id)
}

// Add will add a new customer to the repository
func (mr *CustomerMemoryRepository) Add(cust models.Customer) error {
	if mr.customers == nil {
		mr.Lock()
		mr = &CustomerMemoryRepository{customers: make(map[string]models.Customer)}
		mr.Unlock()
	}
	if _, ok := mr.customers[cust.ID]; ok {
		return errors.InternalError.Newf("fail to add a new customer %v", cust)
	}
	mr.Lock()
	mr.customers[cust.ID] = cust
	mr.Unlock()
	return nil
}

// List customers
func (mr *CustomerMemoryRepository) List() ([]models.Customer, error) {
	cc := make([]models.Customer, 0)
	for _, cust := range mr.customers {
		cc = append(cc, cust)
	}
	return cc, nil
}

// Empty customers
func (mr *CustomerMemoryRepository) Delete() error {
	for k := range mr.customers {
		delete(mr.customers, k)
	}
	return nil
}
