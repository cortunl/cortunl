package wicd

import (
	"github.com/cortunl/cortunl/dbus"
)

func daemon() (conn *dbus.Conn, err error) {
	conn, err = dbus.NewConn(dbus.SystemBus,
		"org.wicd.daemon", "/org/wicd/daemon", "org.wicd.daemon")
	if err != nil {
		return
	}

	return
}

func wired() (conn *dbus.Conn, err error) {
	conn, err = dbus.NewConn(dbus.SystemBus,
		"org.wicd.daemon", "/org/wicd/daemon/wired", "org.wicd.daemon.wired")
	if err != nil {
		return
	}

	return
}

func wireless() (conn *dbus.Conn, err error) {
	conn, err = dbus.NewConn(dbus.SystemBus,
		"org.wicd.daemon", "/org/wicd/daemon/wireless",
		"org.wicd.daemon.wireless")
	if err != nil {
		return
	}

	return
}
