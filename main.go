package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github-uploader/config"
	"github-uploader/router"
	"github-uploader/server"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init (){
	log.Default().SetOutput(os.Stdout)
	log.Default().SetFlags(log.Llongfile | log.Ldate | log.Ltime)
}

func main(){
	var configPath string
	flag.StringVar(&configPath, "config", "app.yaml", "配置文件路径")
	flag.Parse()

	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取配置文件异常: %s\n", err.Error())
	}

	if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(config.App); err != nil {
		log.Fatalf("解析配置文件异常: %s\n", err.Error())
	}

	fmt.Printf("============== Config ===================\n%s\n=========================================\n", config.App)

	// PID
	var pidConfig = config.App.Pid
	if pidConfig != nil {
		var pid = os.Getpid()
		log.Printf("进程ID: %d\n", pid)
		if pidConfig.File != "" {
			if err := os.WriteFile(pidConfig.File, []byte(fmt.Sprintf("%d", pid)), 0x775); err != nil {
				if pidConfig.FatalOnError {
					log.Fatalf("写入PID文件异常: %s\n", err.Error())
				} else {
					log.Printf("写入PID文件异常: %s\n" + err.Error())
				}
			} else {
				defer func() {
					if err := os.Remove(pidConfig.File); err != nil {
						log.Printf("删除PID文件异常: %s\n" + err.Error())
					}
				}()
			}
		}
	}

	// HTTP服务
	go func() {
		log.Println("HTTP服务启动")
		if err := server.Run(router.GetRouter()); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP服务启动异常: %s\n" + err.Error())
		}
	}()

	notify := make(chan os.Signal, 1)
	signal.Notify(notify)
	for {
		sig := <- notify
		log.Printf("收到信号: %s\n", sig.String())
		switch sig {
			case os.Kill, os.Interrupt: {
				func(){
					// HTTP 服务停止
					ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
					defer cancel()
					log.Println("停止HTTP服务器")
					if err := server.Shutdown(ctx); err != nil {
						log.Printf("HTTP服务器停止异常: %s\n", err.Error())
					}
				}()
				log.Println("Bye")
				return
			}
			default:{
			}
		}
	}
}
