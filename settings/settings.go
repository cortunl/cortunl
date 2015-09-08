package settings

var Settings = &SettingsData{}

type SettingsData struct {
	ExternalWired    []string
	ExternalWireless []string
	InternalWired    []string
	InternalWireless []string
	WirelessSsid     string
	WirelessPassword string
}
