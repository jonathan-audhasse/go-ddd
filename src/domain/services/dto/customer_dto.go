package dto

import (
	"goddd/domain/models"
	"log"
)

type CustomerDTO struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email"`
}

func (dto CustomerDTO) ToModel(userid string) (models.Customer, error) {
	cust, err := models.NewCustomer(userid, dto.Name, dto.Email)
	if err != nil {
		// handle error
		log.Printf("fail to transfer %v to customer model: %s\n", dto, err)
		return models.Customer{}, err
	}
	return cust, nil
}
