package controller

import (
	"errors"
	"github-uploader/common"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// ErrorHandle 包装handler，如果handler返回异常，则停止请求链，并且根据error响应异常信息给客户端
func ErrorHandle (handle func (*gin.Context) error) func (*gin.Context) {
	return func(context *gin.Context) {
		err := handle(context)
		if err != nil {
			ErrorResponse(context, err)
			context.Abort()
		}
	}
}

// ResultHandler 返回响应客户端的JSON结果或者异常信息
func ResultHandler (handle func (*gin.Context) (*common.Response, error))func (*gin.Context)  {
	return func(context *gin.Context)  {
		response, err := handle(context)
		if err != nil {
			ErrorResponse(context, err)
			context.Abort()
		} else {
			context.JSON(response.Code.HttpStatus, response)
		}
	}
}

// FailResponse 响应业务异常信息，err 不能为空
func failResponse (ctx *gin.Context, err *common.ServiceError){
	if err.Response != nil {
		ctx.JSON(err.HttpStatus, err.Response)
	} else {
		ctx.Status(err.HttpStatus)
	}
}

// ErrorResponse 异常处理
func ErrorResponse (ctx *gin.Context, err error) {
	switch rawError := err.(type) {
		case common.ServiceError, *common.ServiceError: {
			if err, ok := rawError.(*common.ServiceError) ; ok {
				failResponse(ctx, err)
			}  else if err, ok := rawError.(common.ServiceError) ; ok {
				failResponse(ctx, &err)
			}
		}
		default: {

			var errMessage = err.Error()

			log.Printf("request err: %s\n", errMessage)

			var responseBody = common.FailResponse(common.CodeServerError.SetMessage(errMessage))

			// 详细的异常类型判断
			if os.IsNotExist(err) {											// 文件未找到
				responseBody = common.FailResponse(common.CodeNotFound)
			} else if os.IsPermission(err) {								// 权限不足
				responseBody = common.FailResponse(common.CodeForbidden)
			} else if err.Error() == "http: request body too large" {		// 消息体过大
				responseBody = common.FailResponse(common.CodeRequestEntityTooLarge.SetMessage("请求体大小超出限制"))
			} else if err == io.EOF { 					 // 绑定请求体时的异常，客户端没有发送请求体
				responseBody = common.FailResponse(common.CodeBadRequest.SetMessage("缺少请求体"))
			} else if strings.HasPrefix(errMessage, "json:") {  // JSON解析异常
				responseBody = common.FailResponse(common.CodeBadRequest.SetMessage("非法请求体"))
			} else if errors.Is(err, http.ErrNotMultipart) {
				responseBody = common.FailResponse(common.CodeBadRequest.SetMessage("不是合法的multipart请求"))
			}

			ctx.JSON(responseBody.Code.HttpStatus, responseBody)
		}
	}
}