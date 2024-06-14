package controller

import (
	"goddd/src/domain/models"
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

// func responseValue(word models.Customer) any {
// 	resBody := struct {
// 		Customer *string `json:"word"`
// 	}{nil}
// 	if !word.IsNil() {
// 		resBody.Customer = &word.Value
// 	}
// 	return resBody
// }

func (clt *CustomerCtl) ListCustomers(c *gin.Context) {
	customers, err := clt.service.ListCustomers()
	if err != nil {
		handleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, customers)
}

// func (clt *CustomerCtl) GetCustomer(c *gin.Context) {
// 	// checking input
// 	if !c.Request.URL.Query().Has("prefix") {
// 		// missing prefix
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "missing 'prefix' parameter"})
// 		return
// 	}
// 	prefix := c.Query("prefix")
// 	w, err := clt.service.GetCustomerFrom(prefix)
// 	if err != nil {
// 		handleErr(c, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, responseValue(w))
// }

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
	switch err {
	case models.ErrInvalidCustomer:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error not handled: oops!"})
	}
}
