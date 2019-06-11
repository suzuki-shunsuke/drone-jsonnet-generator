package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/drone-jsonnet-generator/generator/domain"
	"github.com/suzuki-shunsuke/drone-jsonnet-generator/generator/handler"
)

func main() {
	app := cli.NewApp()
	app.Name = "drone-jsonnet-generator"
	app.Version = domain.Version
	app.Author = "suzuki-shunsuke https://github.com/suzuki-shunsuke"
	app.Usage = "convert Drone v0.8x format .drone.yml to v1 format .drone.jsonnet"
	app.Commands = []cli.Command{
		handler.GenCommand,
	}
	app.Run(os.Args)
}
