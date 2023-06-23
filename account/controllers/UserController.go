package controllers

import (
	"account/context"
	"account/models"
	"account/request"
	"account/response"
	"account/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	httpContext "github.com/gorilla/context"
	"gorm.io/gorm"
)

type UserController struct {
}

var userRep *models.UserRepository
var userService services.UserService
var notExist uint

func init() {
	notExist = 0
	userRep = models.UserRep
	userService = services.UserServiceImpl{}
	govalidator.SetFieldsRequiredByDefault(true)
}

func (u UserController) AddUser(c *gin.Context) {
	var userRequest request.Registre
	var user *models.User

	err := c.ShouldBindJSON(&userRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.AddUser{Err: response.MessageBadeRequest})
		return
	}

	ok, err := govalidator.ValidateStruct(userRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, response.AddUser{ErrValidation: err})
		return
	}

	user = user.FromRequest(userRequest)

	status, res := userService.AddUser(user)

	c.JSON(status, res)

}

func (u UserController) DeleteUser(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.DeleteUser{Err: "Invalid id: " + c.Param("id"), IsDeleted: false})
		return
	}
	user := models.User{Model: gorm.Model{ID: uint(id)}}
	_, errDelete := userRep.Delete(&user)

	if errDelete != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.DeleteUser{Err: response.MessageBadeRequest, IsDeleted: false})
		return
	}

	c.JSON(http.StatusOK, response.DeleteUser{IsDeleted: true})

}

func (u UserController) GetAll(c *gin.Context){
	users:=userRep.GetAll()
	c.JSON(http.StatusOK, users)
}

func (u UserController) UpdateUser(c *gin.Context) {
	var userRequest request.Update
	var user *models.User

	err := c.ShouldBindJSON(&userRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.UpdateUser{Err: response.MessageBadeRequest})
		return
	}

	ok, err := govalidator.ValidateStruct(userRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, response.UpdateUser{ErrValidation: err})
		return
	}
	user = user.FromRequest(userRequest)

	user.Model.ID = uint(c.MustGet("user").(*models.User).ID)
	status, result := userService.UpdateUser(user)
	c.JSON(status, result)

}

func (u UserController) UpdatePassword(c *gin.Context) {
	var updatePassword request.UpdatePassword
	var user *models.User
	err := c.Bind(&updatePassword)
	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.UpdatePassword{Err: response.MessageBadeRequest, IsUpdated: false})
		return
	}
	ok, err := govalidator.ValidateStruct(updatePassword)
	if !ok {
		c.JSON(http.StatusBadRequest, response.UpdatePassword{ErrValidation: err, IsUpdated: false})
		return
	}
	user = c.MustGet("user").(*models.User)

	err = userService.UpdatePassword(updatePassword, *user)

	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.UpdatePassword{Err: err.Error(), IsUpdated: false})
		return
	} else {
		c.JSON(http.StatusOK, response.UpdatePassword{IsUpdated: true})
	}

}

func (u UserController) AddRoleToUser(c *gin.Context) {
	var role models.Role
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.AddRoleToUser{Err: "Invalid id: " + c.Param("id")})
		return
	}
	user := models.UserRep.FindById(int(id))
	status, result := userService.AddRoleToUser(user, role)
	c.JSON(status, result)

}

func (u UserController) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.GetUser{Err: "Invalid id: " + c.Param("id")})
		return
	}

	user := userRep.FindById(int(id))

	if user == nil || user.ID == notExist {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusNotFound, response.GetUser{Err: "User id: " + c.Param("id") + " not found "})
		return
	}

	c.JSON(http.StatusOK, response.GetUser{User: *user})

}

func (u UserController) Login(c *gin.Context) {

	var userRequest request.Login
	var user *models.User
	err := c.ShouldBindJSON(&userRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.LoginResponse{Err: response.MessageBadeRequest})
		return
	}

	ok, err := govalidator.ValidateStruct(userRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, response.LoginResponse{ErrValidation: err})
		return
	}

	user = user.FromRequest(userRequest)

	userIsExist, err := userService.ValidateUser(user)

	if !userIsExist {

		c.JSON(http.StatusBadRequest, response.LoginResponse{Err: err.Error()})
		return
	}

	tokenString, err := userService.GenerateToken(*user)

	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.LoginResponse{Err: response.MessageErrInternal})
		return
	}

	response := response.LoginResponse{Token: tokenString}

	c.JSON(http.StatusOK, response)

}

func (u UserController) ValidateToken(c *gin.Context) {

	var token request.TokenValidation
	err := c.Bind(&token)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ValidateToken{Err: response.MessageBadeRequest})
	}
	ok, err := govalidator.ValidateStruct(token)
	if !ok {
		c.JSON(http.StatusBadRequest, response.ValidateToken{ErrValidation: err})
		return
	}
	user, err := userService.ValidateToken(token.Token)
	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ValidateToken{Err: "Token Invalid !!"})
		return
	}

	c.JSON(http.StatusOK, response.ValidateToken{User: user})
}

func (u UserController) Test(w http.ResponseWriter, r *http.Request) {

	user := httpContext.Get(r, "user")
	fmt.Fprintf(w, "Hello World! ")
	fmt.Println(user)
}
