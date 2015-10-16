package netctl

import (
	"github.com/cortunl/cortunl/utils"
)

func Disconnect() (err error) {
	err = utils.Exec("", "netctl", "stop-all")
	if err != nil {
		return
	}

	return
}
