package routes

import (
	//"account/middlewares"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {

	//Router = mux.NewRouter()
	Router = gin.Default()
	// use middleware Logs in all route of the projet
	// var logs middlewares.Logs
	// Router.Use(logs.LogsMiddelware)

}
