package network

type Wireless struct {
	*Network
	Ssid         string
	SecurityType string
	SecurityData interface{}
}
