package main

import (
	"fmt"
	"os"

	"github.com/cfg8er/cfg8er/internal/serve"
	cli "gopkg.in/urfave/cli.v1"
)

var configPath string
var listen string

func main() {
	app := cli.NewApp()
	app.Name = "Cfg8er"
	app.Usage = "git based configuration hosting service"

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Start http service and serve configured repositories",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Value: "config.yml",
					Usage: "Configuration file path",
				},
				cli.StringFlag{
					Name:  "listen, l",
					Value: "127.0.0.1:8080",
					Usage: "IP address and port to listen on",
				},
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: "Enable debug mode",
				},
			},
			Action: serve.Run,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
