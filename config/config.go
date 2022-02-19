package config

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"gopkg.in/yaml.v3"
)

// cli里的配置参数，使用类型类似firewalld
var (
	Listen   = ":9879"
	Host     = ""
	Username = ""
	Password = ""

	// 日志配置
	LogLevel      = "info"
	LogFile       = ""   // 日志文件存放路径,如果为空，则输出到控制台
	LogType       = ""   // 日志类型，支持 txt 和 json ，默认txt
	LogMaxSize    = 100  //单位M
	LogMaxBackups = 3    // 日志文件保留个数
	LogMaxAge     = 365  // 单位天
	LogCompress   = true // 压缩轮转的日志

	// 版本信息
	APPName    = ""
	Maintainer = ""
	APPVersion = ""
	BuildTime  = ""
	GitCommit  = ""
	GOVERSION  = runtime.Version()
	GOOSARCH   = runtime.GOOS + "/" + runtime.GOARCH
	// 其他配置文件
	ConfigPath = ""
)

type Service struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

type MyConfig struct {
	Service Service `json:"service" yaml:"service"`
}

var Config *MyConfig

func LoadConfig(filepath string) {
	if filepath == "" {
		return
	}
	config, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("read config failed, please check the path: %v , err: %v", filepath, err)
	}
	err = yaml.Unmarshal(config, &Config)
	if err != nil {
		fmt.Printf("Unmarshal to struct, err: %v", err)
	}
	fmt.Printf("LoadConfig: %v", Config)
}
