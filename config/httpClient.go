package config

import "time"

type HttpClientConfig struct {
	Timeout time.Duration `yaml:"Timeout"`
}
