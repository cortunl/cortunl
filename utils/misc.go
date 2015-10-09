package utils

import (
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
)

func EnableIpForwarding() (err error) {
	err = Exec("", "sysctl", "-w", "net.ipv4.ip_forward=1")
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "utils: Failed to enable IP forwarding"),
		}
		return
	}

	return
}
