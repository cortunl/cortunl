package netctl

import (
	"fmt"
	"github.com/cortunl/cortunl/network"
	"github.com/cortunl/cortunl/utils"
)

func connectWired(netwk *network.WiredNetwork) (err error) {
	err = Disconnect(netwk.Interface)
	if err != nil {
		return
	}

	data := fmt.Sprintf(conf,
		"ethernet",
		netwk.Interface,
	)

	err = utils.Write(confPathPrefix+netwk.Interface, data)
	if err != nil {
		return
	}

	err = utils.Exec("", "netctl", "start", confNamePrefix+netwk.Interface)
	if err != nil {
		return
	}

	return
}

func connectWireless(netwk *network.WirelessNetwork) (err error) {
	err = Disconnect(netwk.Interface)
	if err != nil {
		return
	}

	data := fmt.Sprintf(conf,
		"wireless",
		netwk.Interface,
	)
	data += fmt.Sprintf("ESSID='%s'\n", netwk.Ssid)

	for key, val := range netwk.Security.Export() {
		data += fmt.Sprintf("%s='%s'\n", key, val)
	}

	err = utils.Write(confPathPrefix+netwk.Interface, data)
	if err != nil {
		return
	}

	err = utils.Exec("", "netctl", "start", confNamePrefix+netwk.Interface)
	if err != nil {
		return
	}

	return
}

func Connect(netIntf interface{}) (err error) {
	switch netwk := netIntf.(type) {
	case *network.WiredNetwork:
		err = connectWired(netwk)
	case *network.WirelessNetwork:
		err = connectWireless(netwk)
	default:
		panic("wicd: Unknown network type")
	}
	if err != nil {
		return
	}

	return
}
