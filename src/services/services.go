package services

import (
	"fmt"
	"goddd/src/domain/repository"
)

type Services struct {
	CustService *CustomerService
}

// Factory
func NewServices(repo *repository.Repository) (*Services, error) {
	custService, err := NewCustomerService(repo.CustomerRepo)
	if err != nil {
		return &Services{}, fmt.Errorf("services error : Customer services not initialize")
	}
	return &Services{CustService: custService}, nil
}
