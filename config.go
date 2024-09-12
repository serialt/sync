package main

import (
	"github.com/google/go-github/v64/github"
)

var (
	// 版本信息
	appVersion bool // 控制是否显示版本
	APPVersion = "v0.0.2"
	BuildTime  = "2006-01-02 15:04:05"
	GitCommit  = "xxxxxxxxxxx"
	ConfigFile = "config.yaml"
	config     *Config
)

type Log struct {
	Level string `yaml:"logLevel"` // 日志级别，支持debug,info,warn,error,panic
	File  string `yaml:"logFile"`  // 日志文件存放路径,如果为空，则输出到控制台
}

type Config struct {
	Log           Log      `yaml:"log"`
	MirrorRoot    string   `yaml:"mirrorRoot"`
	GithubToken   []string `yaml:"githubToken"`
	RandomSleep   int64    `yaml:"randomSleep"`
	GithubRelease []string `yaml:"githubRelease"`
	ExcludeTxt    []string `yaml:"excludeTxt"`
	LastNum       int64    `yaml:"lastNum"`
}

type GithubClient struct {
	Client     *github.Client
	TokenIndex int64
}

type Release struct {
	Owner     string
	Repo      string
	Version   string
	AssetName string
	AssetID   int64
}
