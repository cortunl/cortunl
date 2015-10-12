package bridge

import (
	"github.com/cortunl/cortunl/utils"
	"net"
)

type Bridge struct {
	Bridge     string
	Network    *net.IPNet
	Network6   *net.IPNet
	Interfaces []string
}

func (b *Bridge) Start() (err error) {
	b.Bridge = reserveBridge()

	err = utils.Exec("", "brctl", "addbr", b.Bridge)
	if err != nil {
		return
	}

	for _, intf := range b.Interfaces {
		_ = utils.Exec("", "ip", "addr", "flush", "dev", intf)
		_ = utils.Exec("", "ip", "link", "set", "dev", intf, "down")

		err = utils.Exec("", "ip", "link", "set", "dev", intf, "up")
		if err != nil {
			return
		}

		err = utils.Exec("", "brctl", "addif", b.Bridge, intf)
		if err != nil {
			return
		}
	}

	addr := utils.CopyNetwork(b.Network)
	utils.IncIp(addr.IP)
	broadcast := utils.GetBroadcast(addr)
	addr6 := utils.CopyNetwork(b.Network6)
	utils.IncIp(addr6.IP)

	err = utils.Exec("", "ip", "link", "set", "dev", b.Bridge, "up")
	if err != nil {
		return
	}

	err = utils.Exec("", "ip", "addr", "add", addr.String(),
		"broadcast", broadcast.String(), "dev", b.Bridge)
	if err != nil {
		return
	}

	err = utils.Exec("", "ip", "-6", "addr", "add", addr6.String(),
		"dev", b.Bridge)
	if err != nil {
		return
	}

	return
}

func (b *Bridge) Stop() (err error) {
	if b.Bridge == "" {
		return
	}

	_ = utils.Exec("", "ip", "link", "set", "dev", b.Bridge, "down")
	_ = utils.Exec("", "brctl", "delbr", b.Bridge)

	releaseBridge(b.Bridge)

	return
}
