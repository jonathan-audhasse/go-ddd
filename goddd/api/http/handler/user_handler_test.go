package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"goddd/api/http/handler"
	"goddd/internal/application/dto"
	"goddd/internal/domain/models"
	service_mocks "goddd/tests/mocks/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	reqBody := `{"unknown_key":"john"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	ctrl.CreateUser(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code) // current behavior
	svc.AssertNotCalled(t, "CreateNewUser", mock.Anything, mock.Anything)

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INVALID_INPUT", res["code"])
	assert.Equal(t, "invalid input: json: unknown field \"unknown_key\"", res["message"])
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	ctx := context.Background()
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	reqBody := `{"username":"john","email":"john@test.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	mockErr := gofakeit.Error()
	svc.EXPECT().CreateNewUser(ctx, mock.Anything).Return(models.User{}, mockErr)

	ctrl.CreateUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	svc.AssertExpectations(t)

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INTERNAL_ERROR", res["code"])
	assert.Equal(t, mockErr.Error(), res["message"])
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	reqBody := `{"username":"john","email":"john@test.com"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	expUser := models.User{
		ID:       uuid.New(),
		Username: "john",
		Email:    "john@test.com",
	}

	createUserReq := dto.CreateUserRequest{
		Username: "john",
		Email:    "john@test.com",
	}

	svc.EXPECT().CreateNewUser(mock.Anything, createUserReq).Return(expUser, nil)
	ctrl.CreateUser(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res models.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.Equal(t, expUser.ID, res.ID)
	require.Equal(t, expUser.Username, res.Username)
	require.Equal(t, expUser.Email, res.Email)

	svc.AssertExpectations(t)
}

func TestUserHandler_CreateUsers_InvalidJSON(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/users/bulk", strings.NewReader("{invalid-json"))
	rec := httptest.NewRecorder()

	ctrl.CreateUser(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code) // current behavior
	svc.AssertNotCalled(t, "CreateNewUsers", mock.Anything, mock.Anything)

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INVALID_INPUT", res["code"])
	assert.Equal(t, "invalid input: invalid character 'i' looking for beginning of object key string", res["message"])
}

func TestUserHandler_CreateUsers_ServiceError(t *testing.T) {
	ctx := context.Background()
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	reqBody := `[{"username":"john","email":"john@test.com"}]`
	req := httptest.NewRequest(http.MethodPost, "/users/bulk", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	mockErr := gofakeit.Error()
	svc.EXPECT().CreateNewUsers(ctx, mock.Anything).Return([]models.User{}, mockErr)

	ctrl.CreateUsers(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	svc.AssertExpectations(t)

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INTERNAL_ERROR", res["code"])
	assert.Equal(t, mockErr.Error(), res["message"])
}

func TestUserHandler_CreateUsers_Success(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	reqBody := `[{"username":"john","email":"john@test.com"}]`
	req := httptest.NewRequest(http.MethodPost, "/users/bulk", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	expUser := []models.User{
		{
			ID:       uuid.New(),
			Username: "john",
			Email:    "john@test.com",
		},
	}

	createUserReq := dto.CreateUsersRequest{
		{
			Username: "john",
			Email:    "john@test.com",
		},
	}

	svc.EXPECT().CreateNewUsers(mock.Anything, createUserReq).Return(expUser, nil)
	ctrl.CreateUsers(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res []models.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.Len(t, res, 1)
	resUser := res[0]
	require.Equal(t, expUser[0].ID, resUser.ID)
	require.Equal(t, expUser[0].Username, resUser.Username)
	require.Equal(t, expUser[0].Email, resUser.Email)

	svc.AssertExpectations(t)
}

func TestUserHandler_GetUser_FailedToParseId(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	invalidId := "invalid-id-format"
	req := httptest.NewRequest(http.MethodGet, "/users/"+invalidId, nil)

	// inject chi URL param manually
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", invalidId)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	ctrl.GetUser(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INVALID_INPUT", res["code"])
	assert.Equal(t, handler.ErrInvalidUserId.Error(), res["message"])

	svc.AssertExpectations(t)
}

func TestUserHandler_GetUser_ServiceError(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	id := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)

	// inject chi URL param manually
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	mockErr := gofakeit.Error()
	svc.EXPECT().GetUser(mock.Anything, id).Return(nil, mockErr)
	ctrl.GetUser(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INTERNAL_ERROR", res["code"])
	assert.Equal(t, mockErr.Error(), res["message"])

	svc.AssertExpectations(t)
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	id := uuid.New()

	expUser := models.User{
		ID:       id,
		Username: "john",
		Email:    "john@test.com",
	}

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)

	// inject chi URL param manually
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	svc.EXPECT().GetUser(mock.Anything, id).Return(&expUser, nil)
	ctrl.GetUser(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res models.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.Equal(t, expUser.ID, res.ID)
	require.Equal(t, expUser.Username, res.Username)
	require.Equal(t, expUser.Email, res.Email)

	svc.AssertExpectations(t)
}

func TestUserHandler_ListUsers_InvalidLimit(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users/?limit=invalid-limit-format", nil)
	rec := httptest.NewRecorder()

	ctrl.ListUsers(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INVALID_INPUT", res["code"])
	assert.Equal(t, handler.ErrInvalidLimit.Error(), res["message"])

	svc.AssertExpectations(t)
}

func TestUserHandler_ListUsers_InvalidCursor(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users/?cursor=invalid-cursor-format", nil)
	rec := httptest.NewRecorder()

	ctrl.ListUsers(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INVALID_INPUT", res["code"])
	assert.Equal(t, handler.ErrInvalidCursor.Error(), res["message"])

	svc.AssertExpectations(t)
}

func TestUserHandler_ListUsers_ServiceError(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	rec := httptest.NewRecorder()

	expReq := dto.ListUsersRequest{
		Limit:  handler.LimitDefaulValue,
		Cursor: nil,
	}

	mockErr := gofakeit.Error()
	svc.EXPECT().ListUsers(mock.Anything, expReq).Return(nil, mockErr)
	ctrl.ListUsers(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INTERNAL_ERROR", res["code"])
	assert.Equal(t, mockErr.Error(), res["message"])

	svc.AssertExpectations(t)
}

func TestUserHandler_ListUsers_Success(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := handler.NewUserHandler(svc)

	limit, cursor := uint(10), gofakeit.UUID()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/?limit=%d&cursor=%s", limit, cursor), nil)
	rec := httptest.NewRecorder()

	expReq := dto.ListUsersRequest{
		Limit:  limit,
		Cursor: lo.ToPtr(cursor),
	}

	expUser := models.User{
		Username: "john",
		Email:    "john@test.com",
	}

	expRes := dto.ListUsersResponse{
		Users:      []models.User{expUser},
		Limit:      limit,
		NextCursor: lo.ToPtr(gofakeit.UUID()), // another cursor
	}

	svc.EXPECT().ListUsers(mock.Anything, expReq).Return(&expRes, nil)
	ctrl.ListUsers(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res dto.ListUsersResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.Len(t, res.Users, 1)
	require.Equal(t, expRes, res)

	svc.AssertExpectations(t)
}
