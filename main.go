package main

import (
	"os"
	"github.com/codegangsta/cli"
	"github.com/cortunl/cortunl/constants"
)

func main() {
	app := cli.NewApp()
	app.Name = constants.NameFormated
	app.Commands = []cli.Command{
		{
			Name:  "version",
			Usage: "Print app version",
			Action: func(c *cli.Context) {
				println("0.1.0")
			},
		},
	}
	app.Run(os.Args)
}
