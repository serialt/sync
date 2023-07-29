package main

import (
	"flag"
	"fmt"

	"log/slog"

	"github.com/serialt/sugar/v2"
)

func init() {
	flag.BoolVar(&appVersion, "v", false, "Display build and version messages")
	flag.StringVar(&ConfigFile, "c", "config.yaml", "Config file")
	flag.StringVar(&AesData, "d", "", "加密的明文")
	flag.Parse()

	err := sugar.LoadConfig(ConfigFile, &config)
	if err != nil {
		config = new(Config)
	}
	slog.SetDefault(sugar.New(
		sugar.WithFile(config.Log.File),
		sugar.WithLevel(config.Log.Level),
	))
	config.DecryptConfig()

}

func main() {
	if appVersion {
		slog.Info("app info",
			"APPVersion", APPVersion,
			"BuildTime", BuildTime,
			"GitCommit", GitCommit,
		)
		return
	}
	if len(AesData) > 0 {
		fmt.Printf("Encrypted string: %v\n", Encrypt(AesData, AesKey))
		return
	}
	Run()
}
