package user

import (
	"context"
	"goddd/internal/application/apperror"
	"goddd/internal/application/dto"
	"goddd/internal/application/transaction"
	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// UserService user service interface
type UserService interface {
	// CreateNewUser creates a new user
	CreateNewUser(context.Context, dto.CreateUserRequest) (models.User, error)
	// GetUser finds a user by its id
	GetUser(context.Context, uuid.UUID) (*models.User, error)
	// ListUsers returns a list of users from a given cursor.
	//
	// It return a limited number of users set in the request.
	ListUsers(context.Context, dto.ListUsersRequest) (*dto.ListUsersResponse, error)
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
	logger := log.Ctx(ctx).With().Any("request", req).Logger()

	logger.Info().Msg("creating a new user...")

	var res models.User
	err := s.tm.Do(ctx, func(ctx context.Context) error {
		// create a new user in repo
		user, err := s.repo.Create(ctx, req.ToNewUser())
		if err != nil {
			log.Err(err).Msg("failed to create a new user")
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
		log.Err(err).Msg("failed to retrieve a new user")
		return nil, apperror.ToAppError(err)
	}

	logger.Info().Msg("user found")

	return user, nil
}

// ListUsers returns a list of users from a given cursor.
//
// It return a limited number of users set in the request.
func (s *service) ListUsers(ctx context.Context, req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	logger := log.Ctx(ctx).With().Any("request", req).Logger()

	logger.Info().Msg("listing user...")

	// parse request to repository page request
	page, err := req.ToPage()
	if err != nil {
		return nil, err
	}

	// retrieve  users from repo
	pageRes, err := s.repo.FindPaged(ctx, page)
	if err != nil {
		log.Err(err).Msg("failed to list users from repo")
		return nil, apperror.ToAppError(err)
	}

	logger.Info().Msg("users found")

	// set next page result if defined
	var nextCursor *string
	if pageRes.NextCursor != nil {
		nextCursor = lo.ToPtr(pageRes.NextCursor.String())
	}

	res := dto.ListUsersResponse{
		Users:      pageRes.Users,
		Limit:      page.Limit,
		NextCursor: nextCursor,
	}

	return &res, nil
}
