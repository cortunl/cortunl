package dhcpd

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
	"net"
	"path/filepath"
	"strings"
)

type Dhcpd struct {
	path      string
	Interface string
	Network   *net.IPNet
}

func (d *Dhcpd) writeConf() (err error) {
	d.path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	d.path = filepath.Join(d.path, confName)

	broadcast := utils.GetBroadcast(d.Network).String()

	router := ""
	start := ""
	end := ""
	for ip := range utils.IterNetwork(d.Network) {
		if router == "" {
			router = ip.String()
		} else if start == "" {
			start = ip.String()
		} else {
			end = ip.String()
		}
	}

	data := fmt.Sprintf(conf,
		d.Network.IP.String(),
		broadcast,
		start,
		end,
		net.IP(d.Network.Mask).String(),
		broadcast,
		router,
		strings.Join([]string{"8.8.8.8", "8.8.4.4"}, ", "),
	)

	err = utils.CreateWrite(d.path, data)
	if err != nil {
		return
	}

	return
}
