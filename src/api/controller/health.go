package controller

import (
	"goddd/src/api/httperror"
	"goddd/src/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(serv *services.HealthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := serv.IsHealthy(); err != nil {
			// not healthy
			httperror.HandleErr(c, err)
			return
		}
		// healthy
		c.JSON(http.StatusOK, map[string]string{})
	}
}
