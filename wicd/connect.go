package wicd

import (
	"github.com/cortunl/cortunl/network"
	"github.com/pacur/pacur/utils"
)

func Connect(netIntf interface{}) (err error) {
	switch net := netIntf.(type) {
	case *network.WiredNetwork:

	case *network.WirelessNetwork:
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
	default:
		panic("wicd: Unknown network type")
	}

	return
}
