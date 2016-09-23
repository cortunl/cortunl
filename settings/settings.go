package settings

import (
	"encoding/json"
	"github.com/cortunl/cortunl/errortypes"
	"github.com/cortunl/cortunl/hostapd"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"net"
	"os"
	"time"
)

var Path = "/etc/cortunl.json"

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
					Interface:    "eth0",
					AllTraffic:   true,
					NatInterface: true,
					Networks:     []*net.IPNet{},
				},
			},
			Outputs: []*Output{
				&Output{
					Interface: "eth1",
				},
				&Output{
					Interface: "wlan0",
				},
			},
			Network:          network,
			Network6:         network6,
			WirelessSsid:     "cortunl",
			WirelessPassword: "cortunl",
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

	Settings.loaded = true
}

func Load() (err error) {
	_, err = os.Stat(Path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			err = &errortypes.ReadError{
				errors.Wrap(err, "config: File stat error"),
			}
		}
		return
	}

	file, err := ioutil.ReadFile(Path)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File read error"),
		}
		return
	}

	err = json.Unmarshal(file, Settings)
	if err != nil {
		err = &errortypes.ReadError{
			errors.Wrap(err, "config: File unmarshal error"),
		}
		return
	}

	Settings.loaded = true

	return
}

type SettingsData struct {
	loaded         bool
	Routers        []*Router
	WirelessDriver string
	BlinkDuration  time.Duration
}
