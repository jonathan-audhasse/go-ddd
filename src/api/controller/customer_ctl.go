package controller

import (
	"errors"
	"goddd/src/domain/models"
	"goddd/src/domain/repository"
	"goddd/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerCtl struct {
	service *services.CustomerService
}

func NewCustomerClt(serv *services.CustomerService) (*CustomerCtl, error) {
	return &CustomerCtl{service: serv}, nil
}

func (clt *CustomerCtl) ListCustomers(c *gin.Context) {
	customers, err := clt.service.ListCustomers()
	if err != nil {
		handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (clt *CustomerCtl) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'id' should not be empty"})
		return
	}
	customer, err := clt.service.GetCustomer(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, customer)
}

// func (clt *CustomerCtl) AddCustomer(c *gin.Context) {
// 	// checking request body
// 	var reqBody struct {
// 		Customer string `json:"word" binding:"required"`
// 	}
// 	if err := c.ShouldBindJSON(&reqBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "body should be of format {'word': 'abc'}",
// 		})
// 		return
// 	}

// 	w, err := clt.service.AddCustomer(reqBody.Customer)
// 	if err != nil {
// 		handleErr(c, err)
// 		return
// 	}
// 	fmt.Println(w)
// 	c.JSON(http.StatusCreated, responseValue(w))
// }

// simple error handler
func handleErr(c *gin.Context, err error) {
	if werr := errors.Unwrap(err); werr != nil {
		handleErr(c, werr)
	}
	switch err {
	case models.ErrInvalidCustomer:
		// invalid customer
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
	case repository.ErrNotFound:
		// item not found in repository
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error not handled: oops!"})
	}
}
