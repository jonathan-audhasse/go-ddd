package api

import (
	"goddd/src/api/controller"
	"goddd/src/domain/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(services *services.Services) *gin.Engine {
	r := gin.Default()

	// middlewares
	r.Use(Cors()) //For CORS
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	})

	// controllers
	custCtl, err := controller.NewCustomerClt(services.CustService)
	if err != nil {
		log.Fatal(err)
	}

	// routes
	r.GET("/customers", custCtl.ListCustomers)
	r.GET("/customers/:id", custCtl.GetCustomer)
	r.DELETE("/customers/:id", custCtl.DeleteCustomer)
	r.POST("/customers", custCtl.AddCustomer)

	// healthcheck
	r.GET("/health", controller.Health)

	return r
}
