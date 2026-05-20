package handler_test

import (
	"encoding/json"
	"goddd/api/http/handler"
	service_mocks "goddd/tests/mocks/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler_IsHealthy_ServiceError(t *testing.T) {
	svc := service_mocks.NewMockHealthService(t)
	ctrl := handler.NewHealthHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mockErr := gofakeit.Error()
	svc.EXPECT().IsHealthy(mock.Anything).Return(mockErr)

	ctrl.IsHealthy(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	assert.Equal(t, "INTERNAL_ERROR", res["code"])
	assert.Equal(t, mockErr.Error(), res["message"])

	svc.AssertExpectations(t)
}

func TestHealthHandler_IsHealthy_Success(t *testing.T) {
	svc := service_mocks.NewMockHealthService(t)
	ctrl := handler.NewHealthHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	svc.EXPECT().IsHealthy(mock.Anything).Return(nil)

	ctrl.IsHealthy(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.Equal(t, map[string]string{}, res)

	svc.AssertExpectations(t)
}
