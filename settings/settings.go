package settings

var Settings = &SettingsData{
	ExternalWired: []string{},
	ExternalWireless: []string{
		"wlan0",
	},
}

type SettingsData struct {
	ExternalWired    []string
	ExternalWireless []string
	InternalWired    []string
	InternalWireless []string
	WirelessSsid     string
	WirelessPassword string
}
