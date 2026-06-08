package handler

import (
	"encoding/json"
	"fmt"
	"goddd/api/http/httperror"
	"goddd/internal/application/apperror"
	"goddd/internal/application/dto"
	userservice "goddd/internal/application/user"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var ErrInvalidUserId = fmt.Errorf("%w: user Id must be of type uuid", apperror.ErrInvalid)

// UserHandler holds user handlers
type UserHandler struct {
	service userservice.UserService
}

// NewUserHandler instantiate a new user handler
func NewUserHandler(service userservice.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser calls CreateUser service
func (c *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	// parse json body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Err(err).Any("body", r.Body).Msg("invalid input")
		httperror.Write(w, apperror.ErrInvalid)
		return
	}

	// call service
	user, err := c.service.CreateNewUser(r.Context(), req)
	if err != nil {
		httperror.Write(w, err)
		return
	}

	// set succeed status in header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// set response body
	_ = json.NewEncoder(w).Encode(user)
}

// GetUser calls GetUser service
func (c *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("failed to parse Id to UUID")
		httperror.Write(w, ErrInvalidUserId)
		return
	}

	// call service
	user, err := c.service.GetUser(r.Context(), id)
	if err != nil {
		httperror.Write(w, err)
		return
	}

	// set succeed status in header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}
