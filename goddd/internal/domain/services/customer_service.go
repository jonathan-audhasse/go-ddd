package services

import (
	"goddd/domain/models"
	"goddd/domain/repository"
	"goddd/domain/services/dto"
	"goddd/pkg/errors"
	"log"
)

// Services that define customer's service implementation
type CustomerService struct {
	custumerRepo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	// Create the CustomerService
	return &CustomerService{repo}
}

func (cs *CustomerService) ListCustomers(userid string) ([]models.Customer, error) {
	customers, err := cs.custumerRepo.ListByUser(userid)
	if err != nil {
		log.Println("fail to list customers: ", err)
		return []models.Customer{}, errors.Wrap(err, "fail to list customers")
	}
	return customers, nil
}

func (cs *CustomerService) GetCustomer(id string) (models.Customer, error) {
	customer, err := cs.custumerRepo.Get(id)
	if err != nil {
		log.Printf("fail to get customer id=%s: %s\n", id, err)
		return models.Customer{}, errors.Wrapf(err, "fail to get customer id=%v", id)
	}
	return customer, nil
}

func (cs *CustomerService) DeleteCustomer(id string) error {
	if err := cs.custumerRepo.Delete(id); err != nil {
		log.Printf("fail to delete customer id=%s: %s\n", id, err)
		return errors.Wrapf(err, "fail to delete customer id=%v", id)
	}
	return nil
}

func (cs *CustomerService) AddCustomer(userid string, data dto.CustomerDTO) (models.Customer, error) {
	customer, err := data.ToModel(userid)
	if err != nil {
		log.Printf("fail to add customer %v: %s", data, err)
		return models.Customer{}, errors.Wrap(err, "fail to add customer")
	}
	if err := cs.custumerRepo.Add(customer); err != nil {
		log.Printf("fail to add customer %v: %s", data, err)
		return models.Customer{}, errors.Wrapf(err, "fail to add customer (%v)", data)
	}
	return customer, nil
}
