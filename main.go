package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cortunl/cortunl/constants"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = constants.NameFormated
	app.Version = constants.Version

	app.Commands = []cli.Command{
		{
			Name:  "version",
			Usage: "Print app version",
			Action: func(c *cli.Context) {
				fmt.Printf("v%s\n", constants.Version)
			},
		},
	}

	app.Run(os.Args)
}
