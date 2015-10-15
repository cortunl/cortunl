package settings

var Settings = &SettingsData{
	LocalDomain:   "network",
	ExternalWired: []string{},
	ExternalWireless: []string{
		"wlan0",
	},
	WirelessDriver: "auto",
	DnsServers: []string{
		"8.8.8.8",
		"8.8.4.4",
	},
	DnsServers6: []string{
		"2001:4860:4860::8888",
		"2001:4860:4860::8844",
	},
}

type SettingsData struct {
	LocalDomain      string
	ExternalWired    []string
	ExternalWireless []string
	InternalWired    []string
	InternalWireless []string
	WirelessSsid     string
	WirelessPassword string
	WirelessDriver   string
	DnsServers       []string
	DnsServers6      []string
}
