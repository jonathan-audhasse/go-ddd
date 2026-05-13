package dto

import (
	"goddd/internal/domain/models"
)

// CreateUserRequest DTO for user creation
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
}

// ToNewUser parses user from DTO to models
func (dto CreateUserRequest) ToNewUser() models.NewUser {
	return models.NewUser{
		Username: dto.Username,
		Email:    dto.Email,
	}
}
