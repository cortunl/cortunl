package wicd

import (
	"github.com/cortunl/cortunl/dbus"
)

func daemon() (conn *dbus.Conn, err error) {
	conn, err = dbus.NewConn(dbus.SystemBus,
		"org.wicd.daemon", "/org/wicd/daemon")
	if err != nil {
		return
	}

	return
}
