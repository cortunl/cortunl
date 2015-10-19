package bridge

import (
	"fmt"
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"net"
)

type Bridge struct {
	running  bool
	Bridge   string
	Network  *net.IPNet
	Network6 *net.IPNet
	Outputs  []*settings.Output
}

func (b *Bridge) Init() {
	if b.Bridge != "" {
		panic("bridge: Bridge already init")
	}
	b.Bridge = reserveBridge()
}

func (b *Bridge) Deinit() {
	if b.Bridge == "" {
		return
	}
	releaseBridge(b.Bridge)
	b.Bridge = ""
}

func (b *Bridge) Start() (err error) {
	if b.running {
		panic("bridge: Bridge already started")
	}
	b.running = true

	err = utils.Exec("", "brctl", "addbr", b.Bridge)
	if err != nil {
		return
	}

	for _, iface := range b.Outputs {
		_ = utils.Exec("", "ip", "link", "set", "dev", iface.Interface, "down")
		_ = utils.Exec("", "ip", "addr", "flush", "dev", iface.Interface)

		err = utils.Exec("", "ip", "link", "set", "dev", iface.Interface, "up")
		if err != nil {
			return
		}

		err = utils.Exec("", "brctl", "addif", b.Bridge, iface.Interface)
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
	if !b.running {
		return
	}
	b.running = false

	_ = utils.Exec("", "ip", "link", "set", "dev", b.Bridge, "down")
	_ = utils.Exec("", "brctl", "delbr", b.Bridge)

	return
}
