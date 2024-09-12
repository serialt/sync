package main

import (
	"flag"
	"fmt"

	"log/slog"

	"github.com/serialt/sugar/v2"
)

func main() {
	flag.BoolVar(&appVersion, "v", false, "Display build and version messages")
	flag.StringVar(&ConfigFile, "c", "config.yaml", "Config file")

	flag.Parse()

	err := sugar.LoadConfig(ConfigFile, &config)
	if err != nil {
		fmt.Println("get config failed: ", err)
		config = new(Config)
	}
	slog.SetDefault(sugar.New(
		sugar.WithLevel(config.Log.Level),
		sugar.WithFile(config.Log.File)),
	)

	if appVersion {
		slog.Info("app info",
			"APPVersion", APPVersion,
			"BuildTime", BuildTime,
			"GitCommit", GitCommit,
		)
		return
	}

	Run()
}
