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
	status := StatusCode(err)

	resp := Response{
		Code:    codeFromError(err),
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(resp)
}

// StatusCode map application error to http status code
func StatusCode(err error) int {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		return http.StatusNotFound

	case errors.Is(err, apperror.ErrUnprocessable):
		return http.StatusUnprocessableEntity

	case errors.Is(err, apperror.ErrUnavailable):
		return http.StatusServiceUnavailable

	case errors.Is(err, apperror.ErrInternal):
		return http.StatusInternalServerError

	default:
		return http.StatusInternalServerError
	}
}

// codeFromError returns a error code from error
func codeFromError(err error) string {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		return "NOT_FOUND"

	case errors.Is(err, apperror.ErrUnprocessable):
		return "UNPROCESSABLE"

	case errors.Is(err, apperror.ErrUnavailable):
		return "SERVICE_UNAVAILABLE"

	default:
		return "INTERNAL_ERROR"
	}
}
