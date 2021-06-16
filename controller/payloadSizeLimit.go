package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PayloadSizeLimit (maxPayloadSize int64) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if maxPayloadSize > 0 {
			ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxPayloadSize)
		}
		ctx.Next()
	}
}

