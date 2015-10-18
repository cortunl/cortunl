package router

import (
	"github.com/cortunl/cortunl/bridge"
	"github.com/cortunl/cortunl/dhcp"
	"github.com/cortunl/cortunl/hostapd"
	"github.com/cortunl/cortunl/iptables"
	"github.com/cortunl/cortunl/routes"
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"net"
	"strings"
	"time"
)

type Router struct {
	Settings       *settings.Router
	bridge         string
	bridgeServer   *bridge.Bridge
	routes         *routes.Routes
	dhcpServer     *dhcp.Dhcp
	hostapdServers []*hostapd.Hostapd
	iptables       *iptables.IpTables
}

func (r *Router) Init() {
	r.hostapdServers = []*hostapd.Hostapd{}

	r.bridgeServer = &bridge.Bridge{
		Network:  r.Settings.Network,
		Network6: r.Settings.Network6,
		Outputs:  r.Settings.Outputs,
	}

	r.routes = &routes.Routes{
		Inputs:   r.Settings.Inputs,
		Network:  r.Settings.Network,
		Network6: r.Settings.Network6,
	}

	for _, output := range r.Settings.Outputs {
		if !strings.HasPrefix(output.Interface, "w") {
			continue
		}

		server := &hostapd.Hostapd{
			Driver:    hostapd.AutoDrv,
			Interface: output.Interface,
			Ssid:      r.Settings.WirelessSsid,
			Password:  r.Settings.WirelessPassword,
			Channel:   r.Settings.WirelessChannel,
		}

		r.hostapdServers = append(r.hostapdServers, server)
	}

	r.dhcpServer = &dhcp.Dhcp{
		LocalDomain: r.Settings.LocalDomain,
		DnsServers:  r.Settings.DnsServers,
		DnsServers6: r.Settings.DnsServers6,
		Network:     r.Settings.Network,
		Network6:    r.Settings.Network6,
	}

	r.iptables = &iptables.IpTables{
		Inputs:   r.Settings.Inputs,
		Network:  r.Settings.Network,
		Network6: r.Settings.Network6,
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

	r.routes.Bridge = r.bridge
	err = r.routes.AddRoutes()
	if err != nil {
		return
	}

	for _, hostapdServer := range r.hostapdServers {
		hostapdServer.Bridge = r.bridge
		err = hostapdServer.Start()
		if err != nil {
			return
		}
	}

	time.Sleep(1 * time.Second)

	r.dhcpServer.Bridge = r.bridge
	err = r.dhcpServer.Start()
	if err != nil {
		return
	}

	r.iptables.Bridge = r.bridge
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

	err = r.routes.RemoveRoutes()
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
