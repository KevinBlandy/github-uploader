package config

import (
	"errors"
	"fmt"
	"github-uploader/util"
	"mime/multipart"
	"strings"
)


// HttpConfig http配置
type HttpConfig struct {
	// 业务请求最大请求体
	MaxRequestBody *util.DataSize	`yaml:"MaxRequestBody"`
	// 最大Multipart内存占用
	Multipart *MultipartConfig		`yaml:"Multipart"`
}


// MultipartConfig multipart请求
type MultipartConfig struct {
	// 单个文件最大体积
	MaxFileSize *util.DataSize		`yaml:"MaxFileSize"`
	// 一次性最多上传的文件数量
	MaxFiles	int					`yaml:"MaxFiles"`
	// 最大请求
	MaxRequestSize *util.DataSize	`yaml:"MaxRequestSize"`
	// 最大内存
	MaxMemory *util.DataSize		`yaml:"MaxMemory"`
	// 允许上传的文件类型
	AllowedFileSuffix []string		`yaml:"AllowedFileSuffix"`
}

// AllowedFile 判断上传文件是否合法
func (m MultipartConfig) AllowedFile (files []*multipart.FileHeader) error {

	if len(files) == 0 {
		return errors.New("文件列表不能为空")
	} else if m.MaxFiles > 0 && m.MaxFiles < len(files){
		return fmt.Errorf("一次最多只能上传:%d 个文件", m.MaxFiles)
	}

	for _, file := range files {
		if err := m.AllowedFileName(file.Filename); err != nil {
			return err
		}
		if err := m.AllowedFileSize(file.Size); err != nil {
			return err
		}
	}
	return nil
}


// AllowedFileSize 判断文件大小是否在限制范围内
func (m MultipartConfig) AllowedFileSize (size int64) error {
	if size == 0 {
		return errors.New("不能上传空文件")
	}
	maxByte := m.MaxRequestSize.ToByte()
	if maxByte > 1 && size > maxByte {
		return fmt.Errorf("文件大小不能超过: %d KB", m.MaxFileSize.ToKilobytes())
	}
	return nil
}

// AllowedFileName 判断文件名称是否是允许上传的文件类型
func (m MultipartConfig) AllowedFileName (fileName string) error {

	if len(m.AllowedFileSuffix) == 0 {
		return  nil
	}

	if fileName == "" {
		return errors.New("文件名称不能为空")
	}
	index := strings.LastIndex(fileName, ".")
	if index == -1 {
		return fmt.Errorf("文件名称必选包含后缀: %s", fileName)
	}

	suffix := fileName[index + 1:]

	for _, v := range m.AllowedFileSuffix {
		if strings.EqualFold(v, suffix) {
			return nil
		}
	}
	return fmt.Errorf("系统只允许上传: %s", strings.Join(m.AllowedFileSuffix, ","))
}