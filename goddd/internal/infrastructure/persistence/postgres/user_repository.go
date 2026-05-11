package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"
	"goddd/internal/infrastructure/persistence/postgres/sqlcgen"
	"goddd/internal/infrastructure/persistence/transaction"
)

var (
	ErrFailedToFindUserById = errors.New("failed to find user by its ID")
	ErrFailedToCreateUser   = errors.New("failed to create a user")
)

// userRepository implements domain/user.Repository using sqlc-generated queries.
type userRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository instantiate a new user repository
func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{pool}
}

// retrieve sqlc queries from connection pool.
// Define the queries from the context transaction if existed.
// If not create a new one from the connection
func getQueries(ctx context.Context, pool *pgxpool.Pool) *sqlcgen.Queries {
	if tx, ok := transaction.GetTx(ctx); ok {
		return sqlcgen.New(tx)
	}

	return sqlcgen.New(pool)
}

// FindByID retrieve a user by its Id
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	// retrieve queries
	q := getQueries(ctx, r.pool)

	logger := log.Ctx(ctx).With().Str("id", id.String()).Logger()

	logger.Debug().Msg("retrieving user by id...")

	row, err := q.GetUserByID(ctx, ToPgUUID(id))
	if err != nil {
		logger.Err(err).Msg("failed to retrieve user by id")
		return nil, ErrFailedToFindUserById
	}

	res := toDomain(row)

	logger.Debug().Msg("user found")
	return &res, nil
}

func (r *userRepository) FindPaged(ctx context.Context, p repository.Page) (repository.PagedResult, error) {

	// retrieve queries
	q := getQueries(ctx, r.pool)

	rows, err := q.ListUsersPaged(ctx, &sqlcgen.ListUsersPagedParams{
		Column1: ToPgUUID(p.Cursor),
		Limit:   int32(p.Limit) + 1, // fetch one extra to detect whether a next page exists
	})
	if err != nil {
		return repository.PagedResult{}, fmt.Errorf("userRepository.FindPaged: %w", err)
	}

	hasNext := len(rows) > p.Limit
	if hasNext {
		rows = rows[:p.Limit]
	}

	result := repository.PagedResult{
		Users: make([]models.User, len(rows)),
	}
	for i, row := range rows {
		result.Users[i] = toDomain(row)
	}
	if hasNext {
		last := result.Users[len(result.Users)-1].ID
		result.NextCursor = &last
	}

	return result, nil
}

// Create add a new user to repository
func (r *userRepository) Create(ctx context.Context, newUser models.NewUser) (models.User, error) {
	// retrieve queries
	q := getQueries(ctx, r.pool)

	logger := log.Ctx(ctx).With().Any("newUser", newUser).Logger()

	logger.Debug().Msg("creating a new user to repo...")

	row, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Email:    newUser.Email,
		Username: newUser.Username,
	})
	if err != nil {
		log.Err(err).Msg("failed to add user to repo")
		return models.User{}, ErrFailedToCreateUser
	}
	// Reflect any DB-generated values (e.g. defaults).
	res := toDomain(row)

	logger.Debug().Msg("succeed to add new user to repo")

	return res, nil
}

func (r *userRepository) Update(ctx context.Context, u models.User) (models.User, error) {
	// retrieve queries
	q := getQueries(ctx, r.pool)

	row, err := q.UpdateUser(ctx, &sqlcgen.UpdateUserParams{
		ID:       ToPgUUID(u.ID),
		Email:    u.Email,
		Username: u.Username,
	})
	if err != nil {
		return models.User{}, fmt.Errorf("userRepository.Update: %w", err)
	}
	res := toDomain(row)

	return res, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// retrieve queries
	q := getQueries(ctx, r.pool)

	if err := q.DeleteUser(ctx, ToPgUUID(id)); err != nil {
		return fmt.Errorf("userRepository.Delete: %w", err)
	}
	return nil
}

// toDomain maps a sqlcgen.User (infrastructure model) to models.User.
// Keeping this conversion here means the domain never imports sqlcgen.
func toDomain(u *sqlcgen.User) models.User {
	return models.User{
		ID:        FromPgUUID(u.ID),
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}
}
