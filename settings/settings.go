package settings

import (
	"time"
)

var Settings = &SettingsData{
	LocalDomain:    "network",
	InputWired:     "eth0",
	InputWireless:  "wlan0",
	WirelessDriver: "auto",
	DnsServers: []string{
		"8.8.8.8",
		"8.8.4.4",
	},
	DnsServers6: []string{
		"2001:4860:4860::8888",
		"2001:4860:4860::8844",
	},
	BlinkDuration: 5 * time.Second,
}

type SettingsData struct {
	LocalDomain      string
	InputWired       string
	InputWireless    string
	OutputWired      []string
	OutputWireless   []string
	WirelessSsid     string
	WirelessPassword string
	WirelessDriver   string
	DnsServers       []string
	DnsServers6      []string
	BlinkDuration    time.Duration
}
