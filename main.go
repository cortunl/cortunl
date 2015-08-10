package main

import (
	"github.com/codegangsta/cli"
	"github.com/cortunl/cortunl/constants"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = constants.NameFormated
	app.Version = constants.Version
	app.Run(os.Args)
}
