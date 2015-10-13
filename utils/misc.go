package utils

import (
	"fmt"
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
	"net"
)

func EnableIpForwarding() (err error) {
	err = Exec("", "sysctl", "-w", "net.ipv4.ip_forward=1")
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "utils: Failed to enable IPv4 forwarding"),
		}
		return
	}

	err = Exec("", "sysctl", "-w", "net.ipv6.conf.all.accept_ra=2")
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "utils: Failed to enable IPv6 accept_ra"),
		}
		return
	}

	err = Exec("", "sysctl", "-w", "net.ipv6.conf.default.accept_ra=2")
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "utils: Failed to enable IPv6 accept_ra"),
		}
		return
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		err = &constants.UnknownError{
			errors.Wrap(err, "utils: Failed to get network interfaces"),
		}
		return
	}

	for _, iface := range ifaces {
		_ = Exec("", "sysctl", "-w",
			fmt.Sprintf("net.ipv6.conf.%s.accept_ra=2", iface.Name))
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
