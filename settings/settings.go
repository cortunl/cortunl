package settings

import (
	"github.com/cortunl/cortunl/hostapd"
	"net"
	"time"
)

var Settings = &SettingsData{
	WirelessDriver: "auto",
	BlinkDuration:  5 * time.Second,
}

func init() {
	_, network, _ := net.ParseCIDR("192.168.32.0/24")
	_, network6, _ := net.ParseCIDR("fd32:3032::/64")

	Settings.Routers = []*Router{
		&Router{
			LocalDomain: "network",
			Inputs: []*Input{
				&Input{
					Interface:  "eth0",
					AllTraffic: true,
				},
			},
			Outputs: []*Output{
				&Output{
					Interface: "wlan0",
				},
			},
			Network:          network,
			Network6:         network6,
			WirelessSsid:     "archtest",
			WirelessPassword: "archtest1234",
			WirelessChannel:  hostapd.AutoChan,
			DnsServers: []string{
				"8.8.8.8",
				"8.8.4.4",
			},
			DnsServers6: []string{
				"2001:4860:4860::8888",
				"2001:4860:4860::8844",
			},
		},
	}
}

type SettingsData struct {
	Routers        []*Router
	WirelessDriver string
	BlinkDuration  time.Duration
}
