package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Recover (ctx *gin.Context){
	defer func() {
		if err := recover(); err != nil {
			log.Printf("request panic: %v", err)
			if !ctx.Writer.Written() {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}()
	ctx.Next()
}
