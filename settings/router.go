package settings

import (
	"net"
)

type Input struct {
	Interface        string
	Address          string // TODO Add AddressDhcp and AddressStatic
	WirelessSsid     string
	WirelessPassword string
	AllTraffic       bool
	Networks         []*net.IPNet
	Networks6        []*net.IPNet
}

type Output struct {
	Interface string
}

type Router struct {
	Inputs           []*Input
	Outputs          []*Output
	Network          *net.IPNet
	Network6         *net.IPNet
	WirelessSsid     string
	WirelessPassword string
	WirelessChannel  int
	LocalDomain      string
	DnsServers       []string
	DnsServers6      []string
}
