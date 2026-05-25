package e2e_test

import (
	"bytes"
	"encoding/json"
	"goddd/internal/application/dto"
	"goddd/internal/domain/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	resetDB(t)

	body := dto.CreateUserRequest{
		Username: "john",
		Email:    "john@test.com",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/users",
		bytes.NewReader(b),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	testHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res models.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.NotEqual(t, uuid.Nil, res.ID)
	require.Equal(t, body.Username, res.Username)
	require.Equal(t, body.Email, res.Email)
	require.False(t, res.CreatedAt.IsZero())
}
