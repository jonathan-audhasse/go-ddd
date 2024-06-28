package controller

import (
	"goddd/src/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	ErrorID errors.ErrorID `json:"errorId"`
	Msg     string         `json:"message"`
}

func NewApiError(err error) ApiError {
	return ApiError{ErrorID: errors.GetID(err), Msg: errors.FullError(err)}
}

func (err ApiError) Error() string {
	return err.Msg
}

// simple error handler
func HandleErr(c *gin.Context, err error) {
	apiErr, ok := err.(ApiError)
	if !ok {
		apiErr = NewApiError(err)
	}
	c.JSON(getHttpStatusCode(apiErr.ErrorID), apiErr)
}

// map internal error code to Http status code
func getHttpStatusCode(errID errors.ErrorID) int {
	switch errID {
	case errors.UnAuthorized:
		return http.StatusUnauthorized
	case errors.InvalidField, errors.RequiredFieldMissing:
		return http.StatusBadRequest
	case errors.InvalidFormat, errors.RepoItemAlreadyExist:
		return http.StatusUnprocessableEntity
	case errors.RepoItemNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
