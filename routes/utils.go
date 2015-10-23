package routes

import (
	"fmt"
	"github.com/cortunl/cortunl/constants"
	"github.com/cortunl/cortunl/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"net"
	"strings"
)

var (
	reservedTables set.Set
)

type Route struct {
	Default   bool
	Interface string
	Address   net.IP
	Network   *net.IPNet
	Command   []string
}

func getRoutes(ipv6 bool) (routes []*Route, err error) {
	routes = []*Route{}
	args := []string{"route"}
	if ipv6 {
		args = append([]string{"-6"}, args...)
	}

	output, err := utils.ExecOutput("", "ip", args...)
	if err != nil {
		return
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		route := &Route{
			Command: fields,
		}

		for i, field := range fields {
			switch field {
			case "dev":
				route.Interface = fields[i+1]
			case "expires":
				// TODO
				fields = append(fields[:i], fields[i+2:]...)
			}
		}

		if fields[0] == "default" {
			route.Default = true
			route.Address = net.ParseIP(fields[2])
		} else {
			if strings.Contains(fields[0], "/") {
				_, route.Network, err = net.ParseCIDR(fields[0])
				if err != nil {
					err = &constants.UnknownError{
						errors.Wrapf(err, "Failed to parse network '%s'",
							fields[0]),
					}
					return
				}
			} else {
				route.Address = net.ParseIP(fields[0])
			}
		}

		routes = append(routes, route)
	}

	return
}

func GetRoutes() ([]*Route, error) {
	return getRoutes(false)
}

func GetRoutes6() ([]*Route, error) {
	return getRoutes(true)
}

type table struct {
	Num  int
	Name string
}

func init() {
	reservedTables = set.NewSet()
}

func reserveTable() (tbl *table) {
	for i := 0; ; i++ {
		num := 9700 + i
		if !reservedTables.Contains(num) {
			reservedTables.Add(num)

			tbl = &table{
				Num:  num,
				Name: fmt.Sprintf("cortunl%d", i),
			}
			return
		}
	}

	return
}

func releaseTable(tbl *table) {
	reservedTables.Remove(tbl.Num)
}
