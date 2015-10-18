package routes

import (
	"fmt"
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"net"
	"strings"
)

type Routes struct {
	table    *table
	Input    *settings.Input
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

func (r *Routes) getRoutes() (routes [][]string, err error) {
	routes = [][]string{}

	inputAddr, err := utils.GetInterfaceAddr(r.Input.Interface)
	if err != nil {
		return
	}

	routes = append(routes, []string{
		"table", r.table.Name,
		"default", "via",
		inputAddr.Gateway.String(),
	})

	routes = append(routes, []string{
		"table", r.table.Name,
		inputAddr.Network.String(),
		"dev", r.Input.Interface,
	})

	routes = append(routes, []string{
		"table", r.table.Name,
		inputAddr.Gateway.String(),
		"dev", r.Input.Interface,
	})

	routes = append(routes, []string{
		"table", r.table.Name,
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

	routes, err := r.getRoutes()
	if err != nil {
		return
	}

	for _, args := range routes {
		args = append([]string{"route", "add"}, args...)
		err = utils.Exec("", "ip", args...)
		if err != nil {
			return
		}
	}

	return
}

func (r *Routes) RemoveRoutes() (err error) {
	if r.table == nil {
		return
	}
	defer r.removeTable()

	routes, err := r.getRoutes()
	if err != nil {
		return
	}

	for _, args := range routes {
		args = append([]string{"route", "del"}, args...)
		_ = utils.Exec("", "ip", args...)
	}

	return
}
