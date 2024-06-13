package dto

import (
	"goddd/src/domain/models"
	"log"
)

type CustomerDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (dto CustomerDTO) ToModel() (models.Customer, error) {
	cust, err := models.NewCustomer(dto.Name, dto.Email)
	if err != nil {
		// handle error
		log.Printf("fail to transfer %v to customer model: %s\n", dto, err)
		return models.Customer{}, err
	}
	return cust, nil
}
