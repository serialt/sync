package main

import (
	"github.com/serialt/sync/cmd"
	"github.com/serialt/sync/config"
)

var (
	APPName    = "cli"
	Maintainer = "tserialt@gmail.com"
	APPVersion = "v0.2"
	BuildTime  = "20060102"
	GitCommit  = "ccccccccccccccc"
)

func main() {
	cmd.Run()
}

func init() {
	config.APPName = APPName
	config.Maintainer = Maintainer
	config.APPVersion = APPVersion
	config.BuildTime = BuildTime
	config.GitCommit = GitCommit
}
