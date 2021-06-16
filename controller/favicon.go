package controller

import (
	"github-uploader/resource"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)


var favicon []byte

// Favicon /favicon.ico
func Favicon(ctx *gin.Context) error {
	file, err := resource.FS.Open("public/favicon.ico")
	if err != nil {
		return err
	}
	defer file.Close()

	ctx.Status(http.StatusOK)
	ctx.Header("Content-Type", "image/x-icon")
	ctx.Header("Cache-Control", "public, max-age=3600")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		log.Printf("读取favicon.ico异常: %s\n" + err.Error())
	}
	return nil
}
