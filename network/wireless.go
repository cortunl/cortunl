package network

import (
	"github.com/cortunl/cortunl/security"
)

type WirelessNetwork struct {
	*Network
	Mac      string
	Ssid     string
	Quality  int
	Channel  int
	Security security.Security
}
