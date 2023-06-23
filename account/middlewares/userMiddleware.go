package middlewares

import (

	"net/http"
	"account/context"
	"account/models"
	"account/request"
	"account/services"

	"github.com/gin-gonic/gin"
)

type UserMiddleware struct{}

var userService services.UserService

func init(){
	userService=services.UserServiceImpl{}
}

func (u UserMiddleware) Login() Middleware {

	return func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			if len(c.Request.Header["Token"])==0{
				c.JSON(http.StatusBadRequest, "You are not authenticated !")
				return
			}
			token := request.TokenValidation{Token:c.Request.Header["Token"][0]}
			user, err := userService.ValidateToken(token.Token)
			if err != nil {
				context.WarningLogger.Println(err)
				c.JSON(http.StatusBadRequest, "You are not authenticated !")
				return
			}
			c.Set("user", user)
			next(c)

		}
	}
}

func (u UserMiddleware) HaveRole(role string) Middleware {

	return func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {

			user := c.MustGet("user").(*models.User)
			for _, userRole := range user.Roles {
				if userRole.RoleName == role {
					next(c)
					return
				}
			}

			c.JSON(http.StatusUnauthorized, "Unauthorized !!")

			return

		}
	}
}

func (u UserMiddleware) HavePrivilige(privilige string) Middleware {

	return func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			user := c.MustGet("user").(*models.User)
			havePrivilige := make(chan bool)
			for _, userRole := range user.Roles {

				go searchPriviligeInRole(userRole, privilige, havePrivilige)

			}

			for i := 0; i < len(user.Roles); i++ {
				if <-havePrivilige {
					next(c)
					return
				}
			}

			c.JSON(http.StatusUnauthorized, "Unauthorized !!")
			return

		}
	}
}

func searchPriviligeInRole(role models.Role, privilige string, channel chan bool) {
	for _, userPrivilige := range role.Priviliges {

		if userPrivilige.Privilige == privilige {
			channel <- true
			return
		}
	}
	channel <- false
}
