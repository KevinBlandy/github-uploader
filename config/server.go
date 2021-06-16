package config

import (
	"github-uploader/util"
	"time"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string 					`yaml:"Host"`
	Port int 						`yaml:"Port"`
	ReadTimeout time.Duration		`yaml:"ReadTimeout"`
	WriteTimeout time.Duration		`yaml:"WriteTimeout"`
	ReadHeaderTimeout time.Duration	`yaml:"ReadHeaderTimeout"`
	IdleTimeout time.Duration		`yaml:"IdleTimeout"`
	MaxHeaderSize *util.DataSize	`yaml:"MaxHeaderSize"`

	Compression *CompressionConfig `yaml:"Compression"`

	Ssl *SslConfig					`yaml:"Ssl"`
}

// CompressionConfig 压缩配置
type CompressionConfig struct {

	Enabled bool `yaml:"Enabled"`
}


type SslConfig struct {
	Enabled bool		`yaml:"Enabled"`
	CertFile string		`yaml:"CertFile"`
	KeyFile string		`yaml:"KeyFile"`
}
