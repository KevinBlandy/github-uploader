package controller

import (
	"encoding/base64"
	"fmt"
	"github-uploader/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Auth (ctx *gin.Context){

	if len(config.App.Account) == 0 {
		ctx.Next()
		return
	}

	authorization := ctx.GetHeader("Authorization")

	if authorization == "" {
		unauthorized(ctx)
		return
	}

	_, ok := user(authorization)
	if !ok {
		unauthorized(ctx)
		return
	}

	ctx.Next()
}

func unauthorized(ctx *gin.Context){
	ctx.Header("WWW-Authenticate", "Basic realm=" + strconv.Quote("Authorization Required"))
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

func user (authorization string) (string, bool) {
	for account, pass := range config.App.Account {
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", account, pass)))
		if authorization == auth {
			return account, true
		}
	}
	return "", false
}