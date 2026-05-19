package userservice

import (
	"context"
	"goddd/internal/application/apperror"
	"goddd/internal/application/dto"
	"goddd/internal/application/transaction"
	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	CreateNewUser(context.Context, dto.CreateUserRequest) (models.User, error)
	GetUser(context.Context, uuid.UUID) (*models.User, error)
}

// service for user's services implementation
type service struct {
	repo repository.UserRepository
	tm   transaction.TransactionManager
}

// NewUserService instantiate the UserService object
func NewUserService(repo repository.UserRepository, tm transaction.TransactionManager) UserService {
	return &service{repo, tm}
}

// CreateNewUser creates a new user
func (s *service) CreateNewUser(ctx context.Context, req dto.CreateUserRequest) (models.User, error) {
	logger := log.Ctx(ctx).With().Any("create_user_request", req).Logger()

	logger.Info().Msg("creating a new user...")

	var res models.User
	err := s.tm.Do(ctx, func(ctx context.Context) error {
		// create a new user in repo
		user, err := s.repo.Create(ctx, req.ToNewUser())
		if err != nil {
			return err
		}
		res = user
		return nil
	})
	if err != nil {
		return models.User{}, apperror.ToAppError(err)

	}

	logger.Info().Msg("new user created")

	return res, nil
}

// GetUser finds a user by its id
func (s *service) GetUser(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	logger := log.Ctx(ctx).With().Str("user_id", userId.String()).Logger()

	logger.Info().Msg("retrieving a new user...")

	// retrieve from repo
	user, err := s.repo.FindByID(ctx, userId)
	if err != nil {
		return nil, apperror.ToAppError(err)
	}

	logger.Info().Msg("user found")

	return user, nil
}
