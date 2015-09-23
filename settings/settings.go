package settings

var Settings = &SettingsData{
	ExternalWired: []string{},
	ExternalWireless: []string{
		"wlan0",
	},
	WirelessDriver: "auto",
	DnsServers: []string{
		"8.8.8.8",
		"8.8.4.4",
	},
}

type SettingsData struct {
	ExternalWired    []string
	ExternalWireless []string
	InternalWired    []string
	InternalWireless []string
	WirelessSsid     string
	WirelessPassword string
	WirelessDriver   string
	DnsServers       []string
}
