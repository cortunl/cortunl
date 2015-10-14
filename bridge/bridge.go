package bridge

import (
	"fmt"
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
	if b.Bridge != "" {
		panic("bridge: Bridge already started")
	}
	b.Bridge = reserveBridge()

	err = utils.Exec("", "brctl", "addbr", b.Bridge)
	if err != nil {
		return
	}

	for _, iface := range b.Interfaces {
		_ = utils.Exec("", "ip", "link", "set", "dev", iface, "down")
		_ = utils.Exec("", "ip", "addr", "flush", "dev", iface)

		err = utils.Exec("", "ip", "link", "set", "dev", iface, "up")
		if err != nil {
			return
		}

		err = utils.Exec("", "brctl", "addif", b.Bridge, iface)
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

	err = utils.Exec("", "sysctl", "-w",
		fmt.Sprintf("net.ipv6.conf.%s.autoconf=0", b.Bridge))
	if err != nil {
		return
	}

	err = utils.Exec("", "sysctl", "-w",
		fmt.Sprintf("net.ipv6.conf.%s.accept_ra=0", b.Bridge))
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
