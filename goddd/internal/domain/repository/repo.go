package repository

import (
	"context"
	"errors"
	"goddd/internal/domain/models"

	"github.com/google/uuid"
)

var (
	// Error linked to Health
	ErrFailToPing = errors.New("failed to ping repository")
	// Error linked to User
	ErrUserNotFound          = errors.New("user not found")
	ErrFailedToFindUserById  = errors.New("failed to find user by its ID")
	ErrFailedToCreateUser    = errors.New("failed to create a user")
	ErrUserEmailAlreadyExist = errors.New("failed to create a user: the email already exists")
	ErrFailedToUpdateUser    = errors.New("failed to update user")
	ErrFailedToDeleteUser    = errors.New("failed to delete user")
)

type Repository struct {
	HealthRepo HealthRepository
	UserRepo   UserRepository
}

// HealthRepository for health operation
type HealthRepository interface {
	// Ping tries to ping the repository
	Ping(context.Context) error
}

// Page carries cursor-based pagination parameters.
// Limit caps the number of rows returned (max 100).
// Cursor is the ID of the last item seen — omit on the first page.
type Page struct {
	Limit  int
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
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindPaged(ctx context.Context, p Page) (PagedResult, error)
	Create(ctx context.Context, user models.NewUser) (models.User, error)
	Update(ctx context.Context, user models.User) (models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
