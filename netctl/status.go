package netctl

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

func Status(iface string) (status bool, err error) {
	conn, err := dbus.New()
	if err != nil {
		err = &constants.UnknownError{
			errors.Wrap(err, "netctl: Failed to connect to systemd dbus"),
		}
		return
	}
	defer conn.Close()

	prop, err := conn.GetUnitProperty(
		fmt.Sprintf("netctl@%s.service",
			confNamePrefix+iface), "ActiveState")
	if err != nil {
		err = &constants.UnknownError{
			errors.Wrap(err, "netctl: Failed to get systemd service status"),
		}
		return
	}

	val := strings.Replace(prop.Value.String(), `"`, "", -1)

	if val == "active" {
		status = true
	}

	return
}
