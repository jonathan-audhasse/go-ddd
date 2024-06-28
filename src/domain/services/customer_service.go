package services

import (
	"goddd/src/domain/models"
	"goddd/src/domain/repository"
	"goddd/src/domain/services/dto"
	"goddd/src/pkg/errors"
	"log"
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
		return []models.Customer{}, errors.Wrap(err, "fail to list customers")
	}
	return customers, nil
}

func (cs *CustomerService) GetCustomer(id string) (models.Customer, error) {
	customer, err := cs.repo.Get(id)
	if err != nil {
		log.Printf("fail to get customer id=%s: %s\n", id, err)
		return models.Customer{}, errors.Wrapf(err, "fail to get customer id=%v", id)
	}
	return customer, nil
}

func (cs *CustomerService) DeleteCustomer(id string) error {
	if err := cs.repo.Delete(id); err != nil {
		log.Printf("fail to delete customer id=%s: %s\n", id, err)
		return errors.Wrapf(err, "fail to delete customer id=%v", id)
	}
	return nil
}

func (cs *CustomerService) AddCustomer(data dto.CustomerDTO) (models.Customer, error) {
	customer, err := data.ToModel()
	if err != nil {
		log.Printf("fail to add customer %v: %s", data, err)
		return models.Customer{}, errors.Wrap(err, "fail to add customer")
	}
	if err := cs.repo.Add(customer); err != nil {
		log.Printf("fail to add customer %v: %s", data, err)
		return models.Customer{}, errors.Wrapf(err, "fail to add customer (%v)", data)
	}
	return customer, nil
}
