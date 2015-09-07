package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cortunl/cortunl/constants"
)

var (
	versionCmd = cli.Command{
		Name:   "version",
		Usage:  "Print app version",
		Action: version,
	}
)

func init() {
	commands = append(commands, versionCmd)
}

func version(c *cli.Context) {
	fmt.Printf("v%s\n", constants.Version)
}
