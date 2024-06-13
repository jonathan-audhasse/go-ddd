package services

import (
	"errors"
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/domain/repository"
	"goddd/src/services/dto"
	"log"
)

var (
	// ErrNotFound is returned when a customer is not found.
	ErrNotFound = errors.New("the item was not found in the repository")
	// ErrFailedToAdd is returned when the could not be added to the repository.
	ErrFailedToAdd = errors.New("failed to add the customer to the repository")
	// ErrUpdateCustomer is returned when the customer could not be updated in the repository.
	ErrListCustomers = errors.New("failed to list customer from the repository")
)

// Services that define customer's service implementation
type CustomerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) (*CustomerService, error) {
	// Create the CustomerService
	cs := &CustomerService{repo}
	return cs, nil
}

func (cs *CustomerService) ListCustomers() ([]models.Customer, error) {
	customers, err := cs.repo.List()
	if err != nil {
		log.Println("fail to list customers: ", err)
		return []models.Customer{}, fmt.Errorf("fail to list customers: %w", err)
	}
	return customers, nil
}

func (cs *CustomerService) GetCustomer(id string) (models.Customer, error) {
	customer, err := cs.repo.Get(id)
	if err != nil {
		log.Println("fail to list customers: ", err)
		return models.Customer{}, fmt.Errorf("fail to get customer id=%v: %w", id, err)
	}
	return customer, nil
}

func (cs *CustomerService) AddCustomer(data dto.CustomerDTO) (models.Customer, error) {
	customer, err := data.ToModel()
	if err != nil {
		log.Printf("fail to add customer %v: %s", data, err)
		return models.Customer{}, fmt.Errorf("fail to add customer: %w", err)
	}
	if err := cs.repo.Add(customer); err != nil {
		log.Printf("fail to add customer %v: %s", data, err)
		return models.Customer{}, fmt.Errorf("fail to add customer: %w", err)
	}
	return customer, nil
}
