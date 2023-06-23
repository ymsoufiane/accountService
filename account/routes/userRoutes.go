package routes

import (
	"account/controllers"
	"account/middlewares"
)

var userController controllers.UserController
var userMiddleware middlewares.UserMiddleware

func init() {

	userRoute := Router.Group("/user")

	userRoute.POST("/login", userController.Login)

	userRoute.POST("/add", middlewares.Middlewares(userController.AddUser,
		userMiddleware.Login(),
		userMiddleware.HavePrivilige("add_user"),
	))

	userRoute.POST("/validateToken", userController.ValidateToken)

	userRoute.PUT("/update", middlewares.Middlewares(userController.UpdateUser,
		userMiddleware.Login(),
	))

	userRoute.PUT("/updatePassword", middlewares.Middlewares(userController.UpdatePassword,
		userMiddleware.Login(),
	))

	userRoute.DELETE("/:id", middlewares.Middlewares(userController.DeleteUser,
		userMiddleware.Login(),
		userMiddleware.HavePrivilige("delete_user")))

	userRoute.GET("/:id", middlewares.Middlewares(userController.GetUser,
		userMiddleware.Login(),
		userMiddleware.HavePrivilige("show_user_info"),
	))

	userRoute.POST("/:id/addRole", middlewares.Middlewares(userController.AddRoleToUser,
		userMiddleware.Login(),
		userMiddleware.HavePrivilige("add_role"),
	))
	userRoute.GET("/all", middlewares.Middlewares(userController.GetAll,
		userMiddleware.Login(),
		//userMiddleware.HavePrivilige("add_role"),
	))

}
