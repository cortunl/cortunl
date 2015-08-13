package network

const (
	Wep  SecurityType = "wep"
	Wpa  SecurityType = "wpa"
	Wpa2 SecurityType = "wpa2"
	None SecurityType = "none"
)

type SecurityType string

type WirelessNetwork struct {
	*Network
	Ssid         string
	Quality      int
	Channel      int
	Security     SecurityType
	SecurityData interface{}
}
