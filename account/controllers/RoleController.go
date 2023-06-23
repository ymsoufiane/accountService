package controllers

import (
	"account/models"
	"account/request"
	"account/response"
	"account/services"
	"net/http"
	"strconv"

	"account/context"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
}

var roleRep models.RoleRepository
var roleService services.RoleService

func init() {
	roleRep = *models.RoleRep
	roleService = services.RoleServiceImpl{}
}

func (roleController RoleController) AddRole(c *gin.Context) {
	var roleRequest request.AddRole
	var role *models.Role
	err := c.Bind(&roleRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.AddRole{Err: response.MessageBadeRequest})
		return
	}
	role = role.FromRequest(roleRequest)

	c.JSON(roleService.AddRole(role))

}

func (roleController RoleController) UpdateRole(c *gin.Context) {

	var roleRequest request.AddRole
	var role *models.Role
	err := c.Bind(&roleRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.AddRole{Err:response.MessageBadeRequest})
		return
	}

	role = role.FromRequest(roleRequest)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest,response.AddRole{Err:"Invalid id :"+c.Param("id")+" !!"})
		return
	}
	
	role.Model.ID = uint(id)

	c.JSON(roleService.UpdateRole(role))

}

func (roleController RoleController) DeleteRole(c *gin.Context) {

	var role models.Role

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid id :"+c.Param("id")+" !!")
		return
	}

	role.Model.ID = uint(id)
	_, err = roleRep.Delete(&role)

	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.MessageErrInternal)
		return
	}

	c.JSON(http.StatusOK, "success role deleted")

}
func (roleController RoleController) GetRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.GetRole{Err: "Invalid id :" + c.Param("id") + " !!"})
		return
	}
	role := roleRep.FindById(int(id))
	if role == nil || role.ID == notExist {
		context.WarningLogger.Println(err)
		c.JSON(http.StatusNotFound, response.GetRole{Err: "Role not found !!"})
		return
	}

	c.JSON(http.StatusOK, response.GetRole{Role: role})
}
