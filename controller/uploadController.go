package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github-uploader/common"
	"github-uploader/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)


var HttpClient = http.Client{
	Timeout: time.Second* 10,
}

func GetUpload (ctx *gin.Context)  {
	ctx.HTML(http.StatusOK, "upload/upload.html", gin.H {
		"title": config.App.MetaData.Title,
	})
}

func PostUpload (ctx *gin.Context) (*common.Response, error){
	defer func() {
		_ = ctx.Request.Body.Close()
	}()

	formData, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}

	defer func() {
		_  = formData.RemoveAll()
	}()

	var files = formData.File["file"]


	if err := config.App.Http.Multipart.AllowedFile(files); err != nil {
		return nil, common.NewServiceError(common.CodeBadRequest.SetMessage(err.Error()))
	}

	var retVal []interface{}

	var githubConfig = config.App.Github

	for _, v := range files {
		result := func() *common.Response {
			file, err := v.Open()
			if err != nil {
				return common.FailResponse(common.CodeServerError.SetMessage(fmt.Sprintf("打开文件异常:%s", err.Error())))
			}
			defer func() {
				_ = file.Close()
			}()


			fileContent, err := io.ReadAll(file)
			if err != nil {
				return common.FailResponse(common.CodeServerError.SetMessage(fmt.Sprintf("读取文件异常:%s", err.Error())))
			}

			requestBody, err := json.Marshal(map[string]interface{}{
				"message": "file upload",
				"content": base64.StdEncoding.EncodeToString(fileContent),
			})

			if err != nil {
				return common.FailResponse(common.CodeServerError.SetMessage(fmt.Sprintf("JSON格式化异常:%s", err.Error())))
			}


			randomFileName := strings.ReplaceAll(uuid.New().String(), "-", "")

			now := time.Now()
			year := fmt.Sprintf("%d", now.Year())
			month := fmt.Sprintf("%02d", now.Month())
			day := fmt.Sprintf("%02d", now.Day())

			contentPath := fmt.Sprintf("%s/%s/%s/%s%s", year, month, day, randomFileName, v.Filename[strings.LastIndex(v.Filename, "."):])

			log.Printf("上传文件: %s\n", v.Filename)
			log.Printf("文件大小: %d\n", v.Size)
			log.Printf("文件路径: %s\n", contentPath)

			request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", githubConfig.Owner, githubConfig.Repository, contentPath), bytes.NewReader(requestBody))
			if err != nil {
				return common.FailResponse(common.CodeServerError.SetMessage(fmt.Sprintf("HTTP请求构建异常:%s", err.Error())))
			}

			request.Header.Add("Accept", "application/json")
			request.Header.Add("Content-Type", "application/json; charset=utf-8")
			request.Header.Add("Authorization", fmt.Sprintf("token %s", githubConfig.AccessToken))

			response, err := HttpClient.Do(request)
			if err != nil {
				return common.FailResponse(common.CodeServerError.SetMessage(fmt.Sprintf("HTTP请求异常:%s", err.Error())))
			}
			defer func() {
				_ = response.Body.Close()
			}()

			responseBody, err := io.ReadAll(response.Body)
			if err != nil && err != io.EOF{
				return common.FailResponse(common.CodeServerError.SetMessage(fmt.Sprintf("读取响应体异常:%s", err.Error())))
			}

			if response.StatusCode != http.StatusCreated {
				failResponse := common.FailResponse(common.CodeServerError.SetMessage(fmt.Sprintf("Github响应异常: status=%d", response.StatusCode)))
				failResponse.Data = string(responseBody) // Github响应的JSON信息
				return failResponse
			}

			url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s/%s", githubConfig.Owner, githubConfig.Repository, contentPath)

			log.Printf("访问地址: %s\n", url)

			return common.OkResponse(url)
		}()

		retVal = append(retVal, result)
	}


	return &common.Response{
		Success: true,
		Data:    retVal,
		Code:    common.CodeCreated,
		Message: "ok",
	}, nil
}
