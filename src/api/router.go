package api

import (
	"goddd/api/controller"
	"goddd/api/middleware"
	"goddd/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(services *services.Services) *gin.Engine {
	r := gin.Default()

	// middlewares
	r.Use(middleware.Cors()) //For CORS
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	})

	// healthcheck
	r.GET("/health", controller.Health(services.HealthService))

	// Basic Authorisation
	auth := middleware.BasicAuth(services.UserService)

	// controllers
	custCtl := controller.NewCustomerClt(services.CustomerService)

	// routes
	r.GET("/customers", auth, custCtl.ListCustomers)
	r.GET("/customers/:id", auth, custCtl.GetCustomer)
	r.DELETE("/customers/:id", auth, custCtl.DeleteCustomer)
	r.POST("/customers", auth, custCtl.AddCustomer)

	return r
}
