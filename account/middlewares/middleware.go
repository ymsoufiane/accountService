package middlewares

import (

	"github.com/gin-gonic/gin"
) 


type Middleware func(gin.HandlerFunc) gin.HandlerFunc

func Middlewares(f gin.HandlerFunc, middlewares ...Middleware) gin.HandlerFunc {
	for i := len(middlewares)-1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}
