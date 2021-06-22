package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server *ServerConfig 				`yaml:"Server"`
	MetaData *MetaDataConfig			`yaml:"MetaData"`
	Pid *PidConfig						`yaml:"Pid"`
	HttpClient *HttpClientConfig		`yaml:"HttpClient"`
	Http *HttpConfig					`yaml:"Http"`
	Account map[string] string			`yaml:"Account"`
	Github GithubConfig					`yaml:"Github"`
}

type PidConfig struct {
	File string				`yaml:"File"`
	FatalOnError bool		`yaml:"FatalOnError"`
}


func (a AppConfig) String () string {
	result, err := yaml.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %s\n", err.Error())
	}
	return string(result)
}

