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
	// CreateNewUser creates a new user.
	CreateNewUser(context.Context, dto.CreateUserRequest) (models.User, error)
	// CreateNewUsers creates new users in a batch mode.
	CreateNewUsers(context.Context, dto.CreateUsersRequest) ([]models.User, error)
	// GetUser finds a user by its id.
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

	// validate input
	if err := req.Validate(); err != nil {
		return models.User{}, err
	}

	// parse request to model
	newUser := models.NewUser{
		Username: req.Username,
		Email:    req.Email,
	}

	var res models.User
	err := s.tm.Do(ctx, func(ctx context.Context) error {

		// create a new user in repo
		user, err := s.repo.Create(ctx, newUser)
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

// CreateNewUsers creates new users in a batch mode
func (s *service) CreateNewUsers(ctx context.Context, req dto.CreateUsersRequest) ([]models.User, error) {
	logger := log.Ctx(ctx).With().Any("newUsers", req).Logger()

	logger.Info().Msg("creating new user...")

	if len(req) == 0 {
		// nothing to do
		logger.Info().Msg("nothing to add")
		return []models.User{}, nil
	}

	// validate input
	if err := req.Validate(); err != nil {
		return []models.User{}, err
	}

	// parse request to model
	newUsers := make([]models.NewUser, len(req))
	for i, newUser := range req {
		newUsers[i] = models.NewUser{Email: newUser.Email, Username: newUser.Username}
	}

	var res []models.User
	err := s.tm.Do(ctx, func(ctx context.Context) error {

		// create new users in batch mode
		users, err := s.repo.BulkCreates(ctx, newUsers)
		if err != nil {
			return err
		}
		res = users
		return nil
	})
	if err != nil {
		return []models.User{}, apperror.ToAppError(err)
	}

	logger.Info().Int("total", len(res)).Msg("new users created")

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
