package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
	"os"
)

var commands = []cli.Command{}

func Run() (app *cli.App, err error) {
	app = cli.NewApp()
	app.Name = constants.NameFormated
	app.Version = constants.Version

	app.Commands = commands

	err = app.Run(os.Args)
	if err != nil {
		err = &RunError{
			errors.Wrap(err, "cmd: Run error"),
		}
		return
	}

	return
}
