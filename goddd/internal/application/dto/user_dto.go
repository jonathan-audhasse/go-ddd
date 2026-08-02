package dto

import (
	"fmt"
	"goddd/internal/application/apperror"
	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	ErrMissingUsername = fmt.Errorf("%w: username is required", apperror.ErrInvalid)
	ErrInvalidEmail    = fmt.Errorf("%w: email is not a valid address", apperror.ErrInvalid)
	ErrInvalidCursor   = fmt.Errorf("%w: cursor Id must be of type uuid", apperror.ErrInvalid)
)

// CreateUserRequest DTO for user creation
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type CreateUsersRequest []CreateUserRequest

// UpdateUserRequest DTO for user update
type UpdateUserRequest struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

// ListUsersRequest DTO to list users from DB
type ListUsersRequest struct {
	Limit  uint    `json:"limit"`
	Cursor *string `json:"cursor"`
}

// ListUsersResponse is the paginated envelope returned by GET /users.
type ListUsersResponse struct {
	Users      []models.User `json:"users"`
	Limit      uint          `json:"limit"`
	NextCursor *string       `json:"next_cursor,omitempty"` // null on last page
}

// Validate checks the request invariants before it reaches the domain.
func (dto UpdateUserRequest) Validate() error {
	if dto.Email == "" {
		// nothing to do
		return nil
	}

	if _, err := mail.ParseAddress(dto.Email); err != nil {
		return fmt.Errorf("%w: email='%s'", ErrInvalidEmail, dto.Email)
	}

	return nil
}

// Validate checks the request invariants before it reaches the domain.
func (dto CreateUserRequest) Validate() error {
	if strings.TrimSpace(dto.Username) == "" {
		return ErrMissingUsername
	}

	if _, err := mail.ParseAddress(dto.Email); err != nil {
		return fmt.Errorf("%w: email='%s'", ErrInvalidEmail, dto.Email)
	}

	return nil
}

// Validate checks the request invariants before it reaches the domain.
func (dto CreateUsersRequest) Validate() error {
	for _, req := range dto {
		if strings.TrimSpace(req.Username) == "" {
			return ErrMissingUsername
		}

		if _, err := mail.ParseAddress(req.Email); err != nil {
			return fmt.Errorf("%w: email='%s'", ErrInvalidEmail, req.Email)
		}
	}

	return nil
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
