package network

const (
	Wpa2 SecurityType = "wpa2"
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
