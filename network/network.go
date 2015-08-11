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
	Id   string
	Name string
	Net  net.IPNet
	Addr net.Addr
	Type NetworkType
}
