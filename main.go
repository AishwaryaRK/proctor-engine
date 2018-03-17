package main

import (
	"os"

	"github.com/gojekfarm/proctor-engine/logger"
	"github.com/gojekfarm/proctor-engine/migrations"
	"github.com/gojekfarm/proctor-engine/server"

	"github.com/urfave/cli"
)

func main() {
	proctorEngine := cli.NewApp()
	proctorEngine.Name = "proctor-engine"
	proctorEngine.Usage = "Handle executing jobs and maintaining their configuration"
	proctorEngine.Version = "0.1.0"
	proctorEngine.Commands = []cli.Command{
		{
			Name:        "migrate",
			Description: "Run database migrations for proctor engine",
			Action: func(c *cli.Context) {
				err := migrations.Up()
				if err != nil {
					panic(err.Error())
				}
				logger.Info("Migration successful")
			},
		},
		{
			Name:        "rollback",
			Description: "Rollback database migrations by one step for proctor engine",
			Action: func(c *cli.Context) {
				err := migrations.DownOneStep()
				if err != nil {
					panic(err.Error())
				}
				logger.Info("Rollback successful")
			},
		},
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start server",
			Action: func(c *cli.Context) error {
				return server.Start()
			},
		},
	}

	proctorEngine.Run(os.Args)
}
