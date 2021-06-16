package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)


// App 全局配置
var App = new(AppConfig)

// Init config的初始化，返回异常信息
func Init (configPath string)  error {
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件异常:%w", err)
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(App)

	if err != nil {
		return fmt.Errorf("解析配置文件异常:%w", err)
	}
	return nil
}