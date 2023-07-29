package main

var (
	// 版本信息
	appVersion bool // 控制是否显示版本
	APPVersion = "v0.0.2"
	BuildTime  = "2006-01-02 15:04:05"
	GitCommit  = "xxxxxxxxxxx"
	ConfigFile = "config.yaml"
	config     *Config

	AesKey       = "Serialt.tang@gmail.com_555555555"
	AesData      string
	MyRepo       []REPO
	BreakWall    []string
	Githubclient *GithubClient
)

type Log struct {
	Level string `yaml:"logLevel"` // 日志级别，支持debug,info,warn,error,panic
	File  string `yaml:"logFile"`  // 日志文件存放路径,如果为空，则输出到控制台
}

type Config struct {
	Log           Log      `yaml:"log"`
	MirrorRoot    string   `yaml:"mirrorRoot"`
	GithubToken   string   `yaml:"githubToken"`
	Encrypt       bool     `yaml:"encrypt"`
	Monitor       []string `yaml:"monitor"`
	Terraform     []string `yaml:"terraform"`
	GithubRelease []string `yaml:"githubRelease"`
	ExcludeTxt    []string `yaml:"excludeTxt"`
}

/// 定义model

type REPO struct {
	Owner string
	Repo  string
}

type GithubClient struct {
	Token string
}

type GithubRelease struct {
	Owner              string
	Repo               string
	Version            string
	AssetName          []string
	AssetID            []int
	BrowserDownloadUrl []string
	Path               string
}
