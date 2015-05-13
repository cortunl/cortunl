package main

import (
	"os"
	"github.com/codegangsta/cli"
	"github.com/cortunl/cortunl/constants"
)

func main() {
	app := cli.NewApp()
	app.Name = constants.NameFormated
	app.Version = "0.1.0"
	app.Run(os.Args)
}
