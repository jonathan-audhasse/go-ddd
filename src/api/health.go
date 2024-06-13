package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	var i interface{}
	c.JSON(http.StatusOK, i)
}
