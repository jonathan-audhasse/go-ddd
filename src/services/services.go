package services

import (
	"fmt"
	"goddd/src/domain/repository"
)

type Services struct {
	CustService *CustomerService
}

// Factory
func NewServices(custRepo repository.CustomerRepository) (*Services, error) {
	cService, err := NewCustomerService(custRepo)
	if err != nil {
		return &Services{}, fmt.Errorf("services error : Customer services not initialize")
	}
	return &Services{CustService: cService}, nil
}
