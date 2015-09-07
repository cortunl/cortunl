package wicd

import (
	"net"
	"strconv"
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
	Strength int
}

func GetStatus() (status *Status, err error) {
	status = &Status{}

	conn, err := daemon()

	call, err := conn.Call("GetConnectionStatus")
	if err != nil {
		return
	}

	data := []interface{}{}
	err = call.Bind(&data)
	if err != nil {
		return
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

		strength, _ := strconv.Atoi(info[2])
		status.Strength = strength
	case 3:
		status.State = Connected
		status.Type = Wired
	}

	if status.State == Connected {
		status.Ip = net.ParseIP(info[0])
	}

	return
}
