package routes

import (
	"fmt"
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"net"
	"strings"
)

type Routes struct {
	routes   [][]string
	table    *table
	Inputs   []*settings.Input
	Bridge   string
	Network  *net.IPNet
	Network6 *net.IPNet
}

func (r *Routes) createTable() (err error) {
	data, err := utils.Read(tablesPath)
	if err != nil {
		return
	}

	if !strings.Contains(data, fmt.Sprintf("%s\n", r.table.Name)) {
		if !strings.HasSuffix(data, "\n") {
			data += "\n"
		}

		data += fmt.Sprintf("%d %s\n", r.table.Num, r.table.Name)

		err = utils.Write(tablesPath, data)
		if err != nil {
			return
		}
	}

	err = utils.Exec("", "ip", "route", "flush", "table", r.table.Name)
	if err != nil {
		return
	}

	return
}

func (r *Routes) removeTable() (err error) {
	data, err := utils.Read(tablesPath)
	if err != nil {
		return
	}

	if !strings.Contains(data, fmt.Sprintf("%s\n", r.table.Name)) {
		return
	}

	data = strings.Replace(data,
		fmt.Sprintf("%d %s\n", r.table.Num, r.table.Name), "", -1)

	err = utils.Write(tablesPath, data)
	if err != nil {
		return
	}

	return
}

func (r *Routes) getRoutes() (err error) {
	r.routes = [][]string{}
	hasDefault := false

	for _, input := range r.Inputs {
		inputAddr, e := utils.GetInterfaceAddr(input.Interface)
		if e != nil {
			err = e
			return
		}

		if input.AllTraffic && !hasDefault {
			hasDefault = true
			r.routes = append(r.routes, []string{
				"default", "via",
				inputAddr.Gateway.String(),
			})

			r.routes = append(r.routes, []string{
				inputAddr.Network.String(),
				"dev", input.Interface,
			})

			r.routes = append(r.routes, []string{
				inputAddr.Gateway.String(),
				"dev", input.Interface,
			})
		} else {
			for _, network := range input.Networks {
				r.routes = append(r.routes, []string{
					network.String(),
					"dev", input.Interface,
				})
			}
		}
	}

	if !hasDefault {
		r.routes = append(r.routes, []string{
			"default", "via",
			"0.0.0.0",
		})
	}

	r.routes = append(r.routes, []string{
		r.Network.String(),
		"dev", r.Bridge,
	})

	return
}

func (r *Routes) AddRoutes() (err error) {
	if r.table != nil {
		panic("routes: Routes already added")
	}
	r.table = reserveTable()

	err = r.createTable()
	if err != nil {
		return
	}

	err = r.getRoutes()
	if err != nil {
		return
	}

	for _, args := range r.routes {
		args = append([]string{"route", "add", "table", r.table.Name}, args...)
		err = utils.Exec("", "ip", args...)
		if err != nil {
			return
		}
	}

	err = utils.Exec("", "ip", "rule", "add",
		"from", r.Network.String(), "lookup", r.table.Name)
	if err != nil {
		return
	}

	return
}

func (r *Routes) RemoveRoutes() (err error) {
	if r.table == nil {
		return
	}
	defer r.removeTable()

	for _, args := range r.routes {
		args = append([]string{"route", "del", "table", r.table.Name}, args...)
		_ = utils.Exec("", "ip", args...)
	}

	utils.Exec("", "ip", "rule", "del",
		"from", r.Network.String(), "lookup", r.table.Name)
	if err != nil {
		return
	}

	return
}
