package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const version = "0.0.1"

var cmdList = []cli.Command{}

func main() {
	app := cli.NewApp()
	app.Name = "go-slack-emoji"
	app.Version = version
	app.Usage = "cli for slack tools"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token,t",
			Usage:  "your slack API token",
			EnvVar: "SLACK_TOKEN,TOKEN",
		},
	}
	app.Commands = cmdList
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type subCmd interface {
	Run(*cli.Context, string) error
}

func action(c *cli.Context, sc subCmd) error {
	g := c.GlobalString("token")
	return sc.Run(c, g)
}
