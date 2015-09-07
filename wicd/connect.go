package wicd

import (
	"github.com/cortunl/cortunl/network"
	"github.com/cortunl/cortunl/utils"
)

func connectWired() (err error) {
	conn, err := wired()
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Call("ConnectWired")
	if err != nil {
		return
	}

	return
}

func connectWireless(net *network.WirelessNetwork) (err error) {
	lock.Lock()
	defer lock.Unlock()

	num, e := getNetworkNum(net.Ssid)
	if e != nil {
		err = e
		return
	}

	for key, val := range net.Security.Properties() {
		err = utils.Exec("", "wicd-cli", "--wireless",
			"--network", num, "--network-property", key,
			"--set-to", val)
		if err != nil {
			return
		}
	}

	err = utils.Exec("", "wicd-cli", "--wireless",
		"--network", num, "--connect")
	if err != nil {
		return
	}

	return
}

func Connect(netIntf interface{}) (err error) {
	err = Disconnect()
	if err != nil {
		return
	}

	switch net := netIntf.(type) {
	case *network.WiredNetwork:
		err = connectWired()
	case *network.WirelessNetwork:
		err = connectWireless(net)
	default:
		panic("wicd: Unknown network type")
	}
	if err != nil {
		return
	}

	return
}
