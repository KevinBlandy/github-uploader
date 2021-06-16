package router

import (
	"github-uploader/config"
	"github-uploader/controller"
	"github-uploader/resource"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetRouter() http.Handler {

	// gin.SetMode(gin.DebugMode)

	multipartConfig := config.App.Http.Multipart

	var router = gin.New()
	router.HandleMethodNotAllowed = true
	if multipartConfig.MaxMemory.ToByte() > 0 {
		router.MaxMultipartMemory = multipartConfig.MaxMemory.ToByte()
	}

	// 加载模板引擎
	templates, err := LoadTemplates(resource.FS, "templates")
	if err != nil {
		log.Fatalf("加载模板引擎异常：%s\n", err.Error())
	}
	router.SetHTMLTemplate(templates)

	router.NoRoute(controller.NotFound)
	router.NoMethod(controller.MethodNotAllowed)

	if config.App.Server.Compression.Enabled {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// cors
	router.Use(controller.Cors)

	// 骗骗人
	router.Use(func(ctx *gin.Context) {
		ctx.Header("Server", "Apache/2.4.27 (Win32) OpenSSL/1.0.2l PHP/7.1.8")
		ctx.Header("x-powered-by", "ThinkPHP")
	})

	// panic
	router.Use(controller.Recover)
	router.GET("/ping", controller.Ping)

	router.GET("/", controller.Default)
	router.GET("/favicon.ico", controller.ErrorHandle(controller.Favicon))

	// 嵌入式静态资源
	staticFS, err := fs.Sub(resource.FS, "public")
	if err != nil {
		log.Fatalf("嵌入式静态资源目录映射异常: %s\n", err.Error())
	} else {
		router.StaticFS("/static", http.FS(staticFS))
	}

	// 文件上传接口，限制Multipart的请求大小
	router.GET("/upload", controller.Auth, controller.GetUpload)
	router.POST("/upload", controller.Auth, controller.PayloadSizeLimit(multipartConfig.MaxRequestSize.ToByte()), controller.ResultHandler(controller.PostUpload))

	return router
}

// LoadTemplates 加载模板引擎
// 模板名称，就是目录的相对路径
// 	index/index.html
// 	default.html
// 	foo/bar/index.html
func LoadTemplates(fileSystem fs.FS, root string) (*template.Template, error) {

	templates := template.New("templates")

	// 模板方法
	templates.Funcs(map[string]interface{}{})

	return templates, fs.WalkDir(fileSystem, root, func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() {
			absPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			err = func() error {
				file, err := fileSystem.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				context, err := io.ReadAll(file)
				if err != nil {
					return err
				}
				// 把文件分隔符，统一替换为 “/”
				templates.New(strings.ReplaceAll(absPath, string(os.PathSeparator), "/")).Parse(string(context))
				return nil
			}()
			return err
		}
		return nil
	})
}
