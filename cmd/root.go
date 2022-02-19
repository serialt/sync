package cmd

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/serialt/cli/config"
	"github.com/serialt/cli/pkg"
)

func env(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

var (
	appVersion bool
)

// LogLevel      = "info"
// LogFile       = "cli.log" // 日志文件存放路径
// LogType       = ""
// LogMaxSize    = 100  //单位M
// LogMaxBackups = 3    // 日志文件保留个数
// LogMaxAge     = 365  // 单位天
// LogCompress   = true // 压缩轮转的日志
// OutputConsole = true
func init() {

	flag.BoolVarP(&appVersion, "version", "v", false, "Display build and version msg")
	flag.StringVar(&config.Listen, "listen", env("ESX_LISTEN", config.Listen), "Listen port")
	flag.StringVar(&config.Host, "host", env("ESX_HOST", config.Host), "URL ESX host ")
	flag.StringVarP(&config.Username, "username", "u", env("ESX_USERNAME", config.Username), "User for ESX")
	flag.StringVarP(&config.Password, "password", "p", env("ESX_PASSWORD", config.Password), "Password for ESX")
	flag.StringVar(&config.LogLevel, "logLevel", env("ESX_LOG", config.LogLevel), "Log level must be, debug or info")
	flag.StringVar(&config.LogFile, "logFile", env("LogFile", config.LogFile), "Logfile path")
	flag.StringVar(&config.LogType, "logType", env("LogType", config.LogType), "Log format, txt or json, default txt")
	flag.IntVar(&config.LogMaxSize, "logMaxSize", config.LogMaxSize, "Size  of logfile,M")
	flag.IntVar(&config.LogMaxBackups, "logMaxBackups", config.LogMaxBackups, "Num of rotate log file")
	flag.IntVar(&config.LogMaxAge, "logMaxAge", config.LogMaxAge, "Time for the log file, Day")
	flag.BoolVar(&config.LogCompress, "logCompress", config.LogCompress, "Compress rotated file")
	flag.StringVarP(&config.ConfigPath, "cfgFile", "c", env("CONFIG", config.ConfigPath), "[EXPERIMENTAL] Path to config yaml file that can enable TLS or authentication.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Println("使用说明")
		flag.PrintDefaults()
	}
	flag.ErrHelp = fmt.Errorf("\n\nSome errors have occurred, check and try again !!! ")

	flag.CommandLine.SortFlags = false
	flag.Parse()

	// register global var
	config.LoadConfig(config.ConfigPath)
	pkg.Logger = pkg.NewLogger()
	pkg.Sugar = pkg.NewSugarLogger()
}

func Run() {

	if appVersion {
		fmt.Printf("APPName: %v\n Maintainer: %v\n Version: %v\n BuildTime: %v\n GitCommit: %v\n GoVersion: %v\n OS/Arch: %v\n",
			config.APPName,
			config.Maintainer,
			config.APPVersion,
			config.BuildTime,
			config.GitCommit,
			config.GOVERSION,
			config.GOOSARCH)
		return
	}

	pkg.Sugar.Info("info log")
	pkg.Sugar.Info(config.ConfigPath)

	// pkg.Sugar.Info(config.LogFile)
}
