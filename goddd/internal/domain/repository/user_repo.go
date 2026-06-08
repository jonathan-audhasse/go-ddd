package repository

import (
	"context"
	"errors"
	"goddd/internal/domain/models"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrFailedToFindUserById  = errors.New("failed to find user by its ID")
	ErrFailedToCreateUser    = errors.New("failed to create a user")
	ErrUserEmailAlreadyExist = errors.New("failed to create a user: the email already exists")
	ErrFailedToUpdateUser    = errors.New("failed to update user")
	ErrFailedToDeleteUser    = errors.New("failed to delete user")
	ErrFailedToListUsers     = errors.New("failed to list users")
)

// Page carries cursor-based pagination parameters.
// Limit caps the number of rows returned (max 100).
// Cursor is the ID of the last item seen — omit on the first page.
type Page struct {
	Limit  uint
	Cursor uuid.UUID // zero value = first page
}

// PagedResult wraps a slice of users with the next cursor.
// NextCursor is nil when there are no further pages.
type PagedResult struct {
	Users      []models.User
	NextCursor *uuid.UUID
}

// UserRepository user repository operations
type UserRepository interface {
	FindByID(context.Context, uuid.UUID) (*models.User, error)
	FindPaged(context.Context, Page) (PagedResult, error)
	Create(context.Context, models.NewUser) (models.User, error)
	BulkCreates(context.Context, []models.NewUser) (int64, error)
	Update(context.Context, models.User) (models.User, error)
	Delete(context.Context, uuid.UUID) error
}
