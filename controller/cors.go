package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors(context *gin.Context) {
	origin := context.GetHeader("Origin")
	if origin != "" {
		context.Header("Access-Control-Allow-Origin", origin)
		requestHeader := context.GetHeader("Access-Control-Request-Headers")
		if requestHeader != "" {
			context.Header("Access-Control-Allow-Headers", requestHeader)
		}
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE, PATCH")
		context.Header("Access-Control-Expose-Headers", "*")
		context.Header("Access-Control-Max-Age", "3000")

		if context.Request.Method == http.MethodOptions {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}

	context.Next()
}
