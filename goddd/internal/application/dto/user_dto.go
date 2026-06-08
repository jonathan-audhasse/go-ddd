package dto

import (
	"fmt"
	"goddd/internal/application/apperror"
	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var ErrInvalidCursor = fmt.Errorf("%w: cursor Id must be of type uuid", apperror.ErrInvalid)

// CreateUserRequest DTO for user creation
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
}

// ToNewUser parses the new user from DTO to models
func (dto CreateUserRequest) ToNewUser() models.NewUser {
	return models.NewUser{
		Username: dto.Username,
		Email:    dto.Email,
	}
}

// ListUsersRequest DTO to list users from DB
type ListUsersRequest struct {
	Limit  uint    `json:"limit"`
	Cursor *string `json:"cursor"`
}

// ToNewUser parses the list request from DTO to repository paged
func (dto ListUsersRequest) ToPage() (repository.Page, error) {
	res := repository.Page{
		Limit:  dto.Limit,
		Cursor: uuid.Nil,
	}

	if dto.Cursor == nil {
		// no cursor provided
		return res, nil
	}
	cursor, err := uuid.Parse(*dto.Cursor)
	if err != nil {
		log.Err(err).Str("cursor", *dto.Cursor).Msg("failed to parse cursor Id to UUID")
		return repository.Page{}, ErrInvalidCursor
	}

	// set the cursor
	res.Cursor = cursor

	return res, nil
}

// ListUsersResponse is the paginated envelope returned by GET /users.
type ListUsersResponse struct {
	Users      []models.User `json:"users"`
	Limit      uint          `json:"limit"`
	NextCursor *string       `json:"next_cursor"` // null on last page
}
