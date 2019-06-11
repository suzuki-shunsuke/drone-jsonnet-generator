package handler

import (
	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/go-cliutil"

	"github.com/suzuki-shunsuke/drone-jsonnet-generator/generator/domain"
	"github.com/suzuki-shunsuke/drone-jsonnet-generator/generator/usecase"
)

// GenCommand is the sub command "gen".
var GenCommand = cli.Command{
	Name:   "gen",
	Usage:  "convert Drone v0.8x format .drone.yml to v1 format .drone.jsonnet",
	Action: Gen,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "source, s",
			Usage: "source .drone.yml path",
			Value: ".drone.yml",
		},
		cli.StringFlag{
			Name:  "target, t",
			Usage: "target .drone.jsonnet path",
			Value: ".drone.jsonnet",
		},
		cli.BoolFlag{
			Name:  "stdout",
			Usage: "output generated jsonnet to stdout",
		},
	},
}

// Gen is the sub command "gen".
func Gen(c *cli.Context) error {
	return cliutil.ConvErrToExitError(usecase.Convert(&domain.ConvertArg{
		Source: c.String("source"),
		Target: c.String("target"),
		Stdout: c.Bool("stdout"),
	}))
}
