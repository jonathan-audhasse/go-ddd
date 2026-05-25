package e2e_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	testHandler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var res map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))

	require.Equal(t, map[string]string{}, res)
}
