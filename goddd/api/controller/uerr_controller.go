package controller

import (
	"goddd/api/httperror"
	userservice "goddd/internal/application/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *userservice.UserService
}

func NewUserClt(serv *userservice.UserService) *UserController {
	return &UserController{service: serv}
}

func (clt *UserController) AddUser(c *gin.Context) {
	userid := c.GetString("user-id")
	if userid == "" {
		httperror.HandleErr(c, errors.InternalError.New("unknow user"))
		return
	}
	// checking request body
	var dto dto.UserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		httperror.HandleErr(c, httperror.NewError(errors.RequiredFieldMissing.New("customer body should contain at least a name")))
		return
	}

	customer, err := clt.service.AddUser(userid, dto)
	if err != nil {
		httperror.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, customer)
}

func (clt *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "'id' should not be empty"})
		return
	}
	customer, err := clt.service.GetUser(id)
	if err != nil {
		httperror.HandleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, customer)
}
