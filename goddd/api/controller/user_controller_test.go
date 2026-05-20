package controller_test

import (
	"context"
	"encoding/json"
	"goddd/api/controller"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserController_CreateUser_InvalidJSON(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := controller.NewUserController(svc)

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{invalid-json"))
	rec := httptest.NewRecorder()

	ctrl.CreateUser(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code) // current behavior
	svc.AssertNotCalled(t, "CreateNewUser", mock.Anything, mock.Anything)

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INVALID_INPUT", res["code"])
	assert.Equal(t, "invalid input", res["message"])
}

func TestUserController_CreateUser_ServiceError(t *testing.T) {
	ctx := context.Background()
	svc := service_mocks.NewMockUserService(t)
	ctrl := controller.NewUserController(svc)

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

func TestUserController_CreateUser_Success(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := controller.NewUserController(svc)

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

func TestUserController_GetUser_FailedToParseId(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := controller.NewUserController(svc)

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
	assert.Equal(t, "invalid input", res["message"])

	svc.AssertExpectations(t)
}

func TestUserController_GetUser_ServiceError(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := controller.NewUserController(svc)

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

func TestUserController_GetUser_Success(t *testing.T) {
	svc := service_mocks.NewMockUserService(t)
	ctrl := controller.NewUserController(svc)

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
