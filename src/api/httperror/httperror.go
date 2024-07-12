package httperror

import (
	"goddd/src/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpError struct {
	ErrorID errors.ErrorID `json:"errorId"`
	Msg     string         `json:"message"`
}

func NewError(err error) httpError {
	httpErr, ok := err.(httpError)
	if !ok {
		httpErr = httpError{ErrorID: errors.GetID(err), Msg: errors.FullError(err)}
	}
	return httpErr
}

func (err httpError) Error() string {
	return err.Msg
}

// simple error handler
func HandleErr(c *gin.Context, err error) {
	httpErr := NewError(err)
	c.JSON(httpStatusCode(httpErr.ErrorID), httpErr)
}

// for abortion
func Abort(c *gin.Context, err error) {
	httpErr := NewError(err)
	c.Error(httpErr)
	c.AbortWithStatusJSON(httpStatusCode(httpErr.ErrorID), httpErr)
}

// map internal error code to Http status code
func httpStatusCode(errID errors.ErrorID) int {
	switch errID {
	case errors.UnAuthorized:
		return http.StatusUnauthorized
	case errors.InvalidField:
		return http.StatusBadRequest
	case errors.InvalidFormat, errors.RepoItemAlreadyExist:
		return http.StatusUnprocessableEntity
	case errors.RepoItemNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
