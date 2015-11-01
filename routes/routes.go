package routes

import (
	"bytes"
	"fmt"
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Routes struct {
	runner   sync.WaitGroup
	stop     bool
	routes   [][]string
	routes6  [][]string
	rules    [][]string
	rules6   [][]string
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

	err = utils.Exec("", "ip", "-6", "route", "flush", "table", r.table.Name)
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
	gatewayMtu := 0

	mainRoutes, err := GetRoutes()
	if err != nil {
		return
	}

	mainRoutes6, err := GetRoutes6()
	if err != nil {
		return
	}

	for _, input := range r.Inputs {
		mtu, e := utils.GetInterfaceMtu6(input.Interface)
		if e != nil {
			err = e
			return
		}
		mtuStr := strconv.Itoa(mtu)

		inputAddr, e := utils.GetInterfaceAddr(input.Interface)
		if e != nil {
			err = e
			return
		}

		if input.AllTraffic && !hasDefault {
			hasDefault = true

			cmd := []string{}
			for _, route := range mainRoutes {
				if route.Default &&
					route.Interface == input.Interface &&
					bytes.Equal(route.Address, inputAddr.Gateway) {

					cmd = route.Command
					break
				}
			}
			if len(cmd) == 0 {
				cmd = []string{
					"default", "via",
					inputAddr.Gateway.String(),
					"dev", input.Interface,
				}
			}
			r.routes = append(r.routes, cmd)

			if inputAddr.Gateway6 != nil {
				gatewayMtu = mtu

				cmd := []string{}
				for _, route := range mainRoutes6 {
					if route.Default &&
						route.Interface == input.Interface &&
						bytes.Equal(route.Address, inputAddr.Gateway6) {

						cmd = route.Command
						break
					}
				}
				if len(cmd) == 0 {
					cmd = []string{
						"default", "via",
						inputAddr.Gateway6.String(),
						"dev", input.Interface,
						"mtu", mtuStr,
					}
				}
				r.routes = append(r.routes, cmd)
			}

			cmd = []string{}
			for _, route := range mainRoutes {
				if route.Interface == input.Interface &&
					route.Network != nil &&
					bytes.Equal(route.Network.IP, inputAddr.Network.IP) &&
					bytes.Equal(route.Network.Mask, inputAddr.Network.Mask) {

					cmd = route.Command
					break
				}
			}
			if len(cmd) == 0 {
				cmd = []string{
					inputAddr.Network.String(),
					"dev", input.Interface,
				}
			}
			r.routes = append(r.routes, cmd)
		} else {
			for _, network := range input.Networks {
				if utils.IsIPNet6(network) {
					r.routes6 = append(r.routes6, []string{
						network.String(),
						"dev", input.Interface,
						"mtu", mtuStr,
					})
				} else {
					r.routes = append(r.routes, []string{
						network.String(),
						"dev", input.Interface,
					})
				}
			}
		}
	}

	if gatewayMtu != 0 {
		err = utils.SetInterfaceMtu6(r.Bridge, gatewayMtu)
		if err != nil {
			return
		}
	}

	r.routes = append(r.routes, []string{
		r.Network.String(),
		"dev", r.Bridge,
	})

	r.routes6 = append(r.routes6, []string{
		r.Network6.String(),
		"dev", r.Bridge,
		"mtu", strconv.Itoa(gatewayMtu),
	})

	return
}

func (r *Routes) getRules() (err error) {
	r.rules = [][]string{}
	r.rules6 = [][]string{}

	r.rules = append(r.rules, []string{
		"from", r.Network.String(),
		"lookup", r.table.Name,
		"prio", "1",
	})
	r.rules6 = append(r.rules6, []string{
		"from", r.Network6.String(),
		"lookup", r.table.Name,
		"prio", "1",
	})

	r.rules = append(r.rules, []string{
		"unreachable", "from", r.Network.String(),
		"prio", "2",
	})
	r.rules6 = append(r.rules6, []string{
		"unreachable", "from", r.Network6.String(),
		"prio", "2",
	})

	return
}

func (r *Routes) AddRoutes() (err error) {
	if r.table != nil {
		panic("routes: Routes already added")
	}
	r.table = reserveTable()
	r.runner = sync.WaitGroup{}
	r.stop = false

	err = r.createTable()
	if err != nil {
		return
	}

	err = r.getRoutes()
	if err != nil {
		return
	}

	err = r.getRules()
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

	for _, args := range r.routes6 {
		args = append([]string{"-6", "route", "add",
			"table", r.table.Name}, args...)
		err = utils.Exec("", "ip", args...)
		if err != nil {
			return
		}
	}

	for _, args := range r.rules {
		args = append([]string{"rule", "add"}, args...)
		err = utils.Exec("", "ip", args...)
		if err != nil {
			return
		}
	}

	for _, args := range r.rules6 {
		args = append([]string{"-6", "rule", "add"}, args...)
		err = utils.Exec("", "ip", args...)
		if err != nil {
			return
		}
	}

	r.runner.Add(1)
	go func() {
		Runner:
		for r.stop != true {
			for i := 0; i < 5; i++ {
				time.Sleep(200 * time.Millisecond)
				if r.stop {
					break Runner
				}
			}

			for _, args := range r.routes {
				args = append([]string{"route", "add",
					"table", r.table.Name}, args...)
				_ = utils.ExecSilent("", "ip", args...)
			}

			for _, args := range r.routes6 {
				args = append([]string{"-6", "route", "add",
					"table", r.table.Name}, args...)
				_ = utils.ExecSilent("", "ip", args...)
			}
		}
		r.runner.Done()
	}()

	return
}

func (r *Routes) RemoveRoutes() (err error) {
	if r.table == nil {
		return
	}
	defer r.removeTable()

	r.stop = true
	r.runner.Wait()

	for _, args := range r.routes {
		args = append([]string{"route", "del", "table", r.table.Name}, args...)
		_ = utils.Exec("", "ip", args...)
	}

	for _, args := range r.routes6 {
		args = append([]string{"-6", "route", "del",
			"table", r.table.Name}, args...)
		_ = utils.Exec("", "ip", args...)
	}

	for _, args := range r.rules {
		args = append([]string{"rule", "del"}, args...)
		_ = utils.Exec("", "ip", args...)
	}

	for _, args := range r.rules {
		args = append([]string{"-6", "rule", "del"}, args...)
		_ = utils.Exec("", "ip", args...)
	}

	return
}
