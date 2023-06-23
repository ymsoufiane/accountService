package routes

import (
	"account/controllers"
	"account/middlewares"
)

var roleController controllers.RoleController

func init() {
	roleRoute := Router.Group("/role")

	roleRoute.POST("/add", middlewares.Middlewares(roleController.AddRole,
		userMiddleware.Login(),
		userMiddleware.HavePrivilige("add_role"),
	))

	roleRoute.PUT("/:id/update", middlewares.Middlewares(roleController.UpdateRole,
		userMiddleware.Login(),
		userMiddleware.HavePrivilige("add_role"),
	))

	roleRoute.DELETE("/:id/delete", middlewares.Middlewares(roleController.DeleteRole,
		userMiddleware.Login(),
		userMiddleware.HavePrivilige("add_role"),
	))

	roleRoute.GET("/:id", roleController.GetRole)

}
