package handler

import (
	"encoding/json"
	"fmt"
	"goddd/api/http/httperror"
	"goddd/internal/application/apperror"
	"goddd/internal/application/dto"
	userservice "goddd/internal/application/user"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	LimitDefaulValue = 50
)

var (
	ErrInvalidUserId = fmt.Errorf("%w: user Id must be of type uuid", apperror.ErrInvalid)
	ErrInvalidCursor = fmt.Errorf("%w: cursor must be of type uuid", apperror.ErrInvalid)
	ErrInvalidLimit  = fmt.Errorf("%w: failed to convert string limit parameter to integer", apperror.ErrInvalid)
)

// UserHandler holds user handlers
type UserHandler struct {
	service userservice.UserService
}

// NewUserHandler instantiate a new user handler
func NewUserHandler(service userservice.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUser calls GetUser service
func (c *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	// check user id of type UUID
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

// ListUsers calls ListUsers service
func (c *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limitStr := q.Get("limit")
	cursorStr := q.Get("cursor")

	limit := LimitDefaulValue // set limit to default value
	if limitStr != "" {
		// check limit is integer
		limitValue, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Err(err).Str("limit", limitStr).Msg("failed to convert limit parameter to integer")
			httperror.Write(w, ErrInvalidLimit)
			return
		}
		limit = limitValue
	}

	var cursor *string
	// check cursor is of type UUID
	if cursorStr != "" {
		if _, err := uuid.Parse(cursorStr); err != nil {
			log.Err(err).Str("cursor", cursorStr).Msg("failed to parse cursor to UUID")
			httperror.Write(w, ErrInvalidCursor)
			return
		}
		cursor = &cursorStr
	}

	// call service
	users, err := c.service.ListUsers(r.Context(), dto.ListUsersRequest{Limit: uint(limit), Cursor: cursor})
	if err != nil {
		httperror.Write(w, err)
		return
	}

	// set succeed status in header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

// CreateUser calls CreateNewUser service
func (c *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// parse json body
	if err := dec.Decode(&req); err != nil {
		log.Err(err).Msg("invalid input")
		httperror.Write(w, fmt.Errorf("%w: %s", apperror.ErrInvalid, err))
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

// CreateUsers calls CreateNewUsers service
func (c *UserHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUsersRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// parse json body
	if err := dec.Decode(&req); err != nil {
		log.Err(err).Msg("invalid input")
		httperror.Write(w, fmt.Errorf("%w: %s", apperror.ErrInvalid, err))
		return
	}

	// call service
	users, err := c.service.CreateNewUsers(r.Context(), req)
	if err != nil {
		httperror.Write(w, err)
		return
	}

	// set succeed status in header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// set response body
	_ = json.NewEncoder(w).Encode(users)
}

// UpdateUser calls UpdateUser service
func (c *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	// check user id of type UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Err(err).Str("id", idStr).Msg("failed to parse Id to UUID")
		httperror.Write(w, ErrInvalidUserId)
		return
	}

	var req dto.UpdateUserRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// parse json body
	if err := dec.Decode(&req); err != nil {
		log.Err(err).Msg("invalid input")
		httperror.Write(w, fmt.Errorf("%w: %s", apperror.ErrInvalid, err))
		return
	}

	// add the user id to the request
	req.Id = id

	// call service
	user, err := c.service.UpdateUser(r.Context(), req)
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
