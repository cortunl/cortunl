package utils

import (
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
)

func EnableIpForwarding() (err error) {
	err = Exec("", "sysctl", "-w", "net.ipv4.ip_forward=1")
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "utils: Failed to enable IPv4 forwarding"),
		}
		return
	}
	err = Exec("", "sysctl", "-w", "net.ipv6.conf.all.forwarding=1")
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "utils: Failed to enable IPv6 forwarding"),
		}
		return
	}

	return
}
