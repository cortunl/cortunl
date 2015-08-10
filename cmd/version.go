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

func version(c *cli.Context) {
	fmt.Printf("v%s\n", constants.Version)
}
