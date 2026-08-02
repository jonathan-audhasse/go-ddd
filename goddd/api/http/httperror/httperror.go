package httperror

import (
	"encoding/json"
	"errors"
	"goddd/internal/application/apperror"
	"net/http"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Write error to http writer
func Write(w http.ResponseWriter, err error) {
	status := statusCode(err)

	resp := Response{
		Code:    codeFromError(err),
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(resp)
}

// StatusCode map application error to http status code
func statusCode(err error) int {
	switch {
	case errors.Is(err, apperror.ErrInvalid):
		return http.StatusBadRequest

	case errors.Is(err, apperror.ErrNotFound):
		return http.StatusNotFound

	case errors.Is(err, apperror.ErrConflict):
		return http.StatusConflict

	case errors.Is(err, apperror.ErrUnavailable):
		return http.StatusServiceUnavailable

	default:
		return http.StatusInternalServerError
	}
}

// codeFromError returns a error code from error
func codeFromError(err error) string {
	switch {
	case errors.Is(err, apperror.ErrInvalid):
		return "INVALID_INPUT"

	case errors.Is(err, apperror.ErrNotFound):
		return "NOT_FOUND"

	case errors.Is(err, apperror.ErrConflict):
		return "CONFLICT_ERROR"

	case errors.Is(err, apperror.ErrUnavailable):
		return "SERVICE_UNAVAILABLE"

	default:
		return "INTERNAL_ERROR"
	}
}
