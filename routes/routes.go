package routes

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
	"net"
	"strings"
)

type Routes struct {
	table    *table
	Input    string
	Output   string
	Network  *net.IPNet
	Network6 *net.IPNet
}

func (r *Routes) createTable() (err error) {
	data, err := utils.Read(tablesPath)
	if err != nil {
		return
	}

	if strings.Contains(data, fmt.Sprintf("%s\n", r.table.Name)) {
		return
	}

	if !strings.HasSuffix(data, "\n") {
		data += "\n"
	}

	data += fmt.Sprintf("%d %s\n", r.table.Num, r.table.Name)

	err = utils.Write(tablesPath, data)
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

func (r *Routes) AddRoutes() (err error) {
	if r.table != nil {
		panic("routes: Routes already added")
	}
	r.table = reserveTable()
	return
}

func (r *Routes) RemoveRoutes() (err error) {
	if r.table != nil {
		panic("routes: Routes already added")
	}
	r.table = reserveTable()

	err = r.removeTable()
	if err != nil {
		return
	}

	return
}
