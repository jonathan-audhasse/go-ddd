package memory

import (
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/domain/repository"
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
	return models.Customer{}, fmt.Errorf("customer %v not found %w", id, repository.ErrNotFound)
}

// Add will add a new customer to the repository
func (mr *CustomerMemoryRepository) Add(cust models.Customer) error {
	if mr.customers == nil {
		mr.Lock()
		mr = &CustomerMemoryRepository{customers: make(map[string]models.Customer)}
		mr.Unlock()
	}
	if _, ok := mr.customers[cust.ID]; ok {
		return fmt.Errorf("fail to add a new customer %w", repository.ErrFailedToAdd)
	}
	mr.Lock()
	mr.customers[cust.ID] = cust
	mr.Unlock()
	return nil
}

// List customers
func (mr *CustomerMemoryRepository) List() ([]models.Customer, error) {
	var cc []models.Customer
	for _, cust := range mr.customers {
		cc = append(cc, cust)
	}
	return cc, nil
}
