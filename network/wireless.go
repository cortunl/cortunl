package network

type Wireless struct {
	*Network
	Ssid         string
	Quality      int
	Channel      int
	SecurityType string
	SecurityData interface{}
}
