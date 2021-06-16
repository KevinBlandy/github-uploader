package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping (ctx *gin.Context)  {
	ctx.String(http.StatusOK, "ok")
}
