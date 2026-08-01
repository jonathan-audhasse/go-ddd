package httperror_test

import (
	"encoding/json"
	"fmt"
	"goddd/api/http/httperror"
	"goddd/internal/application/apperror"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	testCases := []struct {
		name          string
		err           error
		expStatusCode int
		expErrCode    string
	}{
		{
			name:          "apperror.ErrInvalid",
			err:           apperror.ErrInvalid,
			expStatusCode: http.StatusBadRequest,
			expErrCode:    "INVALID_INPUT",
		},
		{
			name:          "apperror.ErrNotFound",
			err:           apperror.ErrNotFound,
			expStatusCode: http.StatusNotFound,
			expErrCode:    "NOT_FOUND",
		},
		{
			name:          "apperror.ErrUnavailable",
			err:           apperror.ErrUnavailable,
			expStatusCode: http.StatusServiceUnavailable,
			expErrCode:    "SERVICE_UNAVAILABLE",
		},
		{
			name:          "apperror.ErrConflict",
			err:           apperror.ErrConflict,
			expStatusCode: http.StatusConflict,
			expErrCode:    "CONFLICT_ERROR",
		},
		{
			name:          "random error",
			err:           gofakeit.Error(),
			expStatusCode: http.StatusInternalServerError,
			expErrCode:    "INTERNAL_ERROR",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			err := fmt.Errorf("%w: mock test", tc.err)

			httperror.Write(rec, err)

			require.Equal(t, tc.expStatusCode, rec.Code)
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var resp httperror.Response
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

			assert.Equal(t, tc.expErrCode, resp.Code)
			assert.Equal(t, err.Error(), resp.Message)
		})
	}
}
