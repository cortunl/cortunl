package netctl

import (
	"github.com/cortunl/cortunl/utils"
)

func Disconnect(iface string) (err error) {
	err = utils.Exec("", "netctl", "stop", confNamePrefix+iface)
	if err != nil {
		return
	}

	_ = utils.Exec("", "ip", "link", "set", iface, "down")

	return
}

func DisconnectAll() (err error) {
	err = utils.Exec("", "netctl", "stop-all")
	if err != nil {
		return
	}

	return
}
