package netctl

import (
	"fmt"
	"github.com/cortunl/cortunl/network"
	"github.com/cortunl/cortunl/utils"
	"sync"
)

var (
	lock = sync.Mutex{}
)

func connectWireless(netwk *network.WirelessNetwork) (err error) {
	_ = utils.Exec("", "ip", "link", "set", netwk.Interface, "down")

	lock.Lock()
	defer lock.Unlock()

	data := fmt.Sprintf(conf,
		"wireless",
		netwk.Interface,
	)
	data += fmt.Sprintf("ESSID='%s'\n", netwk.Ssid)

	for key, val := range netwk.Security.Properties() {
		data += fmt.Sprintf("%s='%s'\n", key, val)
	}

	err = utils.Write(confPath, data)
	if err != nil {
		return
	}

	err = utils.Exec("", "netctl", "start", confName)
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

	switch netwk := netIntf.(type) {
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
