package userservice

import (
	"context"
	"goddd/internal/application/dto"
	"goddd/internal/application/transaction"
	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// UserService for user's services implementation
type UserService struct {
	repo repository.UserRepository
	tm   transaction.TransactionManager
}

// NewUserService instantiate the UserService object
func NewUserService(repo repository.UserRepository, tm transaction.TransactionManager) *UserService {
	return &UserService{repo, tm}
}

// CreateNewUser creates a new user
func (us *UserService) CreateNewUser(ctx context.Context, req dto.CreateUserRequest) (models.User, error) {
	logger := log.Ctx(ctx).With().Any("create_user_request", req).Logger()

	logger.Info().Msg("creating a new user...")

	// create a new user in repo
	user, err := us.repo.Create(ctx, req.ToNewUser())
	if err != nil {
		return models.User{}, err
	}

	logger.Info().Msg("new user created")

	return user, nil
}

// GetUser finds a user by its id
func (us *UserService) GetUser(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	logger := log.Ctx(ctx).With().Str("user_id", userId.String()).Logger()

	logger.Info().Msg("retrieving a new user...")

	// retrieve from repo
	user, err := us.repo.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("user found")

	return user, nil
}
