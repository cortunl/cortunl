package network

type Wireless struct {
	*Network
	Ssid           string
	EncryptionType string
	EncryptionData interface{}
}
