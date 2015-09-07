package wicd

import (
	"github.com/guelfey/go.dbus"
	"net"
)

const (
	Connected    ConnectionState = "connected"
	Disconnected ConnectionState = "disconnected"
	Connecting   ConnectionState = "connecting"
	Wired        ConnectionType  = "wired"
	Wireless     ConnectionType  = "wireless"
	Unknown      ConnectionType  = "unknown"
)

type ConnectionState string
type ConnectionType string

type Status struct {
	State    ConnectionState
	Type     ConnectionType
	Ip       net.IP
	Ssid     string
	Strength string
}

func GetStatus() (status *Status, err error) {
	status = &Status{}

	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}

	daemon := conn.Object("org.wicd.daemon", "/org/wicd/daemon")
	call := daemon.Call("GetConnectionStatus", 0)
	if call.Err != nil {
		panic(call.Err)
	}

	data := []interface{}{}
	err = call.Store(&data)
	if err != nil {
		panic(err)
	}

	state := int(data[0].(uint32))
	info := data[1].([]string)

	switch state {
	case 0:
		status.State = Disconnected
		status.Type = Unknown
	case 1:
		status.State = Connecting
		switch info[0] {
		case "wired":
			status.Type = Wired
		case "wireless":
			status.Type = Wireless
		default:
			status.Type = Unknown
		}
	case 2:
		status.State = Connected
		status.Type = Wireless
		status.Ssid = info[1]
		status.Strength = info[2]
	case 3:
		status.State = Connected
		status.Type = Wired
	}

	if status.State == Connected {
		status.Ip = net.ParseIP(info[0])
	}

	return
}
