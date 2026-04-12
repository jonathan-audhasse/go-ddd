package repository

import (
	"context"
	"goddd/internal/domain/models"

	"github.com/google/uuid"
)
 
type Repository struct {
	UserRepo     UserRepository
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
 
// Repository is the port the domain defines. The infrastructure layer
// provides the concrete implementation — the domain never imports postgres.
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindPaged(ctx context.Context, p Page) (PagedResult, error)
	Create(ctx context.Context, u *models.User) error
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
 