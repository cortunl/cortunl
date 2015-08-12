package network

type WirelessNetwork struct {
	*Network
	Ssid         string
	Quality      int
	Channel      int
	SecurityType string
	SecurityData interface{}
}
