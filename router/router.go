package router

import (
	"fmt"
	"github.com/cortunl/cortunl/bridge"
	"github.com/cortunl/cortunl/dhcp"
	"github.com/cortunl/cortunl/hostapd"
	"github.com/cortunl/cortunl/iptables"
	"github.com/cortunl/cortunl/routes"
	"github.com/cortunl/cortunl/runner"
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"strings"
	"time"
)

type Router struct {
	Settings *settings.Router
	stopped  bool
	errors   chan error
	brdg     *bridge.Bridge
	routes   *routes.Routes
	dcp      *dhcp.Dhcp
	hstpd    []*hostapd.Hostapd
	iptables *iptables.IpTables
}

func (r *Router) onError(err error) {
	if !r.stopped {
		r.errors <- err
	}
}

func (r *Router) Conf() {
	r.hstpd = []*hostapd.Hostapd{}

	r.brdg = &bridge.Bridge{
		Network:  r.Settings.Network,
		Network6: r.Settings.Network6,
		Outputs:  r.Settings.Outputs,
	}
	r.brdg.Init()

	r.routes = &routes.Routes{
		Inputs:   r.Settings.Inputs,
		Bridge:   r.brdg.Bridge,
		Network:  r.Settings.Network,
		Network6: r.Settings.Network6,
	}

	for _, output := range r.Settings.Outputs {
		if !strings.HasPrefix(output.Interface, "w") {
			continue
		}

		server := &hostapd.Hostapd{
			Runner: runner.Runner{
				OnError: r.onError,
			},
			Driver:    hostapd.AutoDrv,
			Interface: output.Interface,
			Bridge:    r.brdg.Bridge,
			Ssid:      r.Settings.WirelessSsid,
			Password:  r.Settings.WirelessPassword,
			Channel:   r.Settings.WirelessChannel,
		}

		r.hstpd = append(r.hstpd, server)
	}

	r.dcp = &dhcp.Dhcp{
		Runner: runner.Runner{
			OnError: r.onError,
		},
		Bridge:      r.brdg.Bridge,
		LocalDomain: r.Settings.LocalDomain,
		DnsServers:  r.Settings.DnsServers,
		DnsServers6: r.Settings.DnsServers6,
		Network:     r.Settings.Network,
		Network6:    r.Settings.Network6,
	}

	r.iptables = &iptables.IpTables{
		Inputs:   r.Settings.Inputs,
		Bridge:   r.brdg.Bridge,
		Network:  r.Settings.Network,
		Network6: r.Settings.Network6,
	}
}

func (r *Router) Start() (err error) {
	r.Conf()
	r.stopped = false

	for _, input := range r.Settings.Inputs {
		if strings.HasPrefix(input.Interface, "w") {
			netwks, e := netctl.GetNetworks(input.Interface)
			if e != nil {
				err = e
				return
			}

			var netwk *network.WirelessNetwork

			for _, n := range netwks {
				if n.Ssid == input.WirelessSsid {
					netwk = n
					break
				}
			}

			if netwk == nil {
				err = &ConnectError{
					errors.Newf("router: Failed to find wireless network "+
						"'%s' on '%s'", input.WirelessSsid, input.Interface),
				}
				return
			}

			netwk.Security.Set("password", input.WirelessPassword)

			err = netctl.Connect(netwk)
			if err != nil {
				return
			}
		} else {
			netwk := &network.WiredNetwork{
				Network: &network.Network{
					Interface: input.Interface,
				},
			}

			err = netctl.Connect(netwk)
			if err != nil {
				return
			}
		}
	}

	err = utils.EnableIpForwarding()
	if err != nil {
		return
	}

	err = r.brdg.Start()
	if err != nil {
		return
	}

	err = r.routes.AddRoutes()
	if err != nil {
		return
	}

	for _, hostapdServer := range r.hstpd {
		err = hostapdServer.Start()
		if err != nil {
			return
		}
	}

	time.Sleep(1 * time.Second)

	err = r.dcp.Start()
	if err != nil {
		return
	}

	r.iptables.Init()
	err = r.iptables.AddRules()
	if err != nil {
		return
	}

	return
}

func (r *Router) Stop() (err error) {
	err = r.brdg.Stop()
	if err != nil {
		return
	}
	r.brdg.Deinit()

	err = r.routes.RemoveRoutes()
	if err != nil {
		return
	}

	for _, hostapdServer := range r.hstpd {
		err = hostapdServer.Stop()
		if err != nil {
			return
		}

		hostapdServer.Wait()
	}

	err = r.dcp.Stop()
	if err != nil {
		return
	}
	r.dcp.Wait()

	err = r.iptables.RemoveRules()
	if err != nil {
		return
	}

	for _, input := range r.Settings.Inputs {
		netctl.Disconnect(input.Interface)
	}

	return
}
