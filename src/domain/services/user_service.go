package services

import (
	"goddd/domain/models"
	"goddd/domain/repository"
	"goddd/pkg/errors"
)

// UserService for user's services implementation
type UserService struct {
	repo repository.UserRepository
}

// Create the UserService
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo}
}

// Find a user by its username
func (us *UserService) GetUserByUsername(username string) (models.User, error) {
	u, err := us.repo.GetByUsername(username)
	if err != nil {
		return models.User{}, errors.Wrap(err, "failed to retrieve user")
	}
	return u, nil
}
