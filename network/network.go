package network

import (
	"net"
)

const (
	Wired    NetworkType = "wired"
	Wireless NetworkType = "wireless"
)

type NetworkType string

type Network struct {
	Id        string
	Name      string
	Interface string
	Network   *net.IPNet
	Address   net.IP
	Network6  *net.IPNet
	Address6  net.IP
	Type      NetworkType
}
