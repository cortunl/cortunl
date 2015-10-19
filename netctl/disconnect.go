package netctl

import (
	"github.com/cortunl/cortunl/utils"
)

func Disconnect(iface string) (err error) {
	err = utils.Exec("", "netctl", confNamePrefix+iface)
	if err != nil {
		return
	}

	return
}
