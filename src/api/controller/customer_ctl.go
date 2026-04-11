package controller

import (
	"goddd/api/httperror"
	"goddd/domain/services"
	"goddd/domain/services/dto"
	"goddd/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerCtl struct {
	service *services.CustomerService
}

func NewCustomerClt(serv *services.CustomerService) *CustomerCtl {
	return &CustomerCtl{service: serv}
}

func (clt *CustomerCtl) ListCustomers(c *gin.Context) {
	userid := c.GetString("user-id")
	if userid == "" {
		httperror.HandleErr(c, errors.InternalError.New("unknow user"))
		return
	}
	customers, err := clt.service.ListCustomers(userid)
	if err != nil {
		httperror.HandleErr(c, err)
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
		httperror.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (clt *CustomerCtl) AddCustomer(c *gin.Context) {
	userid := c.GetString("user-id")
	if userid == "" {
		httperror.HandleErr(c, errors.InternalError.New("unknow user"))
		return
	}
	// checking request body
	var dto dto.CustomerDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		httperror.HandleErr(c, httperror.NewError(errors.RequiredFieldMissing.New("customer body should contain at least a name")))
		return
	}

	customer, err := clt.service.AddCustomer(userid, dto)
	if err != nil {
		httperror.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, customer)
}

func (clt *CustomerCtl) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'id' should not be empty"})
		return
	}
	if err := clt.service.DeleteCustomer(id); err != nil {
		httperror.HandleErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
