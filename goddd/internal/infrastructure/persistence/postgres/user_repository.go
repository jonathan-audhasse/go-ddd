package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"
	"goddd/internal/infrastructure/persistence/postgres/sqlcgen"
)

// UserRepository implements domain/user.Repository using sqlc-generated queries.
type UserRepository struct {
	q *sqlcgen.Queries
}

func NewUserRepository(tx pgx.Tx) *UserRepository {
	return &UserRepository{q: sqlcgen.New(tx)}
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.FindByID: %w", err)
	}
	return toDomain(row), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.FindByEmail: %w", err)
	}
	return toDomain(row), nil
}

func (r *UserRepository) FindPaged(ctx context.Context, p repository.Page) (repository.PagedResult, error) {
	rows, err := r.q.ListUsersPaged(ctx, &sqlcgen.ListUsersPagedParams{
		Cursor: p.Cursor,
		Limit:  int32(p.Limit) + 1, // fetch one extra to detect whether a next page exists
	})
	if err != nil {
		return repository.PagedResult{}, fmt.Errorf("UserRepository.FindPaged: %w", err)
	}

	hasNext := len(rows) > p.Limit
	if hasNext {
		rows = rows[:p.Limit]
	}

	result := repository.PagedResult{
		Users: make([]models.User, len(rows)),
	}
	for i, row := range rows {
		result.Users[i] = *toDomain(row)
	}
	if hasNext {
		last := result.Users[len(result.Users)-1].ID
		result.NextCursor = &last
	}

	return result, nil
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	row, err := r.q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("UserRepository.Create: %w", err)
	}
	// Reflect any DB-generated values (e.g. defaults) back into the entity.
	*u = *toDomain(row)
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	row, err := r.q.UpdateUser(ctx, &sqlcgen.UpdateUserParams{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
	})
	if err != nil {
		return fmt.Errorf("UserRepository.Update: %w", err)
	}
	*u = *toDomain(row)
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("UserRepository.Delete: %w", err)
	}
	return nil
}

// toDomain maps a sqlcgen.User (infrastructure model) to models.User.
// Keeping this conversion here means the domain never imports sqlcgen.
func toDomain(u *sqlcgen.User) *models.User {
	return &models.User{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}