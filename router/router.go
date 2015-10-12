package router

import (
	"github.com/cortunl/cortunl/bridge"
	"github.com/cortunl/cortunl/dhcp"
	"github.com/cortunl/cortunl/hostapd"
	"github.com/cortunl/cortunl/iptables"
	"net"
	"strings"
	"github.com/cortunl/cortunl/utils"
)

type Router struct {
	Inputs         []string
	Outputs        []string
	Network        *net.IPNet
	Network6       *net.IPNet
	Ssid           string
	Password       string
	bridge         string
	bridgeServer   *bridge.Bridge
	dhcpServer     *dhcp.Dhcp
	hostapdServers []*hostapd.Hostapd
	iptables       *iptables.IpTables
}

func (r *Router) Init() {
	r.hostapdServers = []*hostapd.Hostapd{}

	r.bridgeServer = &bridge.Bridge{
		Network:    r.Network,
		Network6:   r.Network6,
		Interfaces: r.Outputs,
	}

	for _, output := range r.Outputs {
		if !strings.HasPrefix(output, "w") {
			continue
		}

		server := &hostapd.Hostapd{
			Driver:    hostapd.Realtek, // TODO
			Interface: output,
			Ssid:      r.Ssid,
			Channel:   1, // TODO
			Password:  r.Password,
		}

		r.hostapdServers = append(r.hostapdServers, server)
	}

	r.dhcpServer = &dhcp.Dhcp{
		Network:  r.Network,
		Network6: r.Network6,
	}

	r.iptables = &iptables.IpTables{
		Input: r.Inputs[0], // TODO
		Network:  r.Network,
		Network6: r.Network6,
	}
}

func (r *Router) Start() (err error) {
	err = utils.EnableIpForwarding()
	if err != nil {
		return
	}

	err = r.bridgeServer.Start()
	if err != nil {
		return
	}
	r.bridge = r.bridgeServer.Bridge

	for _, hostapdServer := range r.hostapdServers {
		err = hostapdServer.Start()
		if err != nil {
			return
		}
	}

	r.dhcpServer.Interface = r.bridge
	r.dhcpServer.Init()
	err = r.dhcpServer.Start()
	if err != nil {
		return
	}

	r.iptables.Output = r.bridge
	r.iptables.Init()
	err = r.iptables.AddRules()
	if err != nil {
		return
	}

	return
}

func (r *Router) Stop() (err error) {
	err = r.bridgeServer.Stop()
	if err != nil {
		return
	}

	for _, hostapdServer := range r.hostapdServers {
		err = hostapdServer.Stop()
		if err != nil {
			return
		}

		hostapdServer.Wait()
	}

	err = r.dhcpServer.Stop()
	if err != nil {
		return
	}
	r.dhcpServer.Wait()

	err = r.iptables.RemoveRules()
	if err != nil {
		return
	}

	return
}
