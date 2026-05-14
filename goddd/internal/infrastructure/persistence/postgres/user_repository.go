package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"
	"goddd/internal/infrastructure/persistence/postgres/sqlcgen"
	"goddd/internal/infrastructure/persistence/transaction"
)

const usersEmailKey = "users_email_key"

var (
	ErrFailedToBulkInsert = errors.New("failed to bulk insert users")
	ErrFailedToRetrieveTx = errors.New("failed to retrieve transaction")
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
	logger := log.Ctx(ctx).With().Str("id", id.String()).Logger()

	logger.Debug().Msg("retrieving user by id...")

	// retrieve queries
	q := getQueries(ctx, r.pool)

	row, err := q.GetUserByID(ctx, ToPgUUID(id))
	if err != nil && err.Error() == noRowErrMessage {
		logger.Err(err).Msg("user not found")
		return nil, repository.ErrUserNotFound
	}
	if err != nil {
		logger.Err(err).Msg("failed to retrieve user by id")
		return nil, repository.ErrFailedToFindUserById
	}

	res := toDomain(row)

	logger.Debug().Msg("user found")
	return &res, nil
}

// FindPaged list users from a page pointed by a cursor.
// It returns a limit of user and the next cursor for the following queries
func (r *userRepository) FindPaged(ctx context.Context, page repository.Page) (repository.PagedResult, error) {

	logger := log.Ctx(ctx).With().
		Str("cursor", page.Cursor.String()).
		Int("limit", page.Limit).
		Logger()

	logger.Debug().Msg("listing users...")

	// retrieve queries
	q := getQueries(ctx, r.pool)

	rows, err := q.ListUsersPaged(ctx, &sqlcgen.ListUsersPagedParams{
		Column1: ToPgUUID(page.Cursor),
		Limit:   int32(page.Limit) + 1, // fetch one extra to detect whether a next page exists
	})
	if err != nil {
		logger.Err(err).Msg("failed to list users")
		return repository.PagedResult{}, repository.ErrFailedToListUsers
	}

	hasNext := len(rows) > page.Limit
	if hasNext {
		rows = rows[:page.Limit]
	}

	res := repository.PagedResult{
		Users: make([]models.User, len(rows)),
	}
	for i, row := range rows {
		res.Users[i] = toDomain(row)
	}
	if hasNext {
		last := res.Users[len(res.Users)-1].ID
		res.NextCursor = &last
	}

	logger.Debug().Int("total", len(res.Users)).Msg("users found")
	return res, nil
}

// Create adds a new user to repository
func (r *userRepository) Create(ctx context.Context, newUser models.NewUser) (models.User, error) {
	logger := log.Ctx(ctx).With().Any("newUser", newUser).Logger()

	logger.Debug().Msg("creating a new user to repo...")

	// retrieve queries
	q := getQueries(ctx, r.pool)

	row, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Email:    newUser.Email,
		Username: newUser.Username,
	})
	if err != nil {
		log.Err(err).Msg("failed to add user to repo")

		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			// email already exists
			if pqErr.Code == uniqueConstraintViolation && pqErr.ConstraintName == usersEmailKey {
				return models.User{}, repository.ErrUserEmailAlreadyExist
			}
		}

		return models.User{}, repository.ErrFailedToCreateUser
	}
	// Reflect any DB-generated values (e.g. defaults).
	res := toDomain(row)

	logger.Debug().Msg("succeed to add new user to repo")

	return res, nil
}

// BulkCreates bulk insert a list of users in the repository
func (*userRepository) BulkCreates(ctx context.Context, users []models.NewUser) (int64, error) {
	logger := log.Ctx(ctx).With().Int("total", len(users)).Logger()

	logger.Debug().Msg("bulk insert user in repo...")

	if len(users) == 0 {
		return 0, nil
	}

	// retrieve transaction
	tx, ok := transaction.GetTx(ctx)
	if !ok {
		return 0, ErrFailedToRetrieveTx
	}

	// append rows
	rows := make([][]any, 0, len(users))
	for _, u := range users {
		rows = append(rows, []any{u.Email, u.Username})
	}

	// bulk insert
	count, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"email", "username"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Err(err).Msg("failed to add users in repo")
		return 0, ErrFailedToBulkInsert
	}

	logger.Debug().Int64("count", count).Msg("users added")

	return count, nil
}

// Update updates the user data
func (r *userRepository) Update(ctx context.Context, user models.User) (models.User, error) {
	logger := log.Ctx(ctx).With().Any("user", user).Logger()

	logger.Debug().Msg("updating user repo data...")

	// retrieve queries
	q := getQueries(ctx, r.pool)

	row, err := q.UpdateUser(ctx, &sqlcgen.UpdateUserParams{
		ID:       ToPgUUID(user.ID),
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		log.Err(err).Msg("failed to update user repo")

		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			// email already exists
			if pqErr.Code == uniqueConstraintViolation && pqErr.ConstraintName == usersEmailKey {
				return models.User{}, repository.ErrUserEmailAlreadyExist
			}
		}

		return models.User{}, repository.ErrFailedToUpdateUser
	}
	res := toDomain(row)

	return res, nil
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	logger := log.Ctx(ctx).With().Str("user_id", id.String()).Logger()

	logger.Debug().Msg("deleting user from repo...")

	// retrieve queries
	q := getQueries(ctx, r.pool)

	if err := q.DeleteUser(ctx, ToPgUUID(id)); err != nil {
		logger.Err(err).Msg("failed to delete user")
		return repository.ErrFailedToDeleteUser
	}
	return nil
}

// toDomain maps a sqlcgen.User (infrastructure model) to models.User.
// Keeping this conversion here means the domain never imports sqlcgen.
func toDomain(u *sqlcgen.User) models.User {
	return models.User{
		ID:        u.ID.Bytes,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}
}
