package network

import (
	"net"
)

type Network struct {
	Id   string
	Name string
	Net  net.IPNet
	Addr net.Addr
	Type string
}
