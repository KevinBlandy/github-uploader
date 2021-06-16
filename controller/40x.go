package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFound (ctx *gin.Context){
	ctx.AbortWithStatus(http.StatusNotFound)
}

func MethodNotAllowed(ctx *gin.Context){
	ctx.AbortWithStatus(http.StatusMethodNotAllowed)
}
