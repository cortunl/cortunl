package utils

import (
	"bytes"
	"container/list"
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"net"
	"strings"
)

type InterfaceAddr struct {
	Gateway  net.IP
	Address  net.IP
	Network  *net.IPNet
	Gateway6 net.IP
	Address6 net.IP
	Network6 *net.IPNet
}

var offsets = [...]int{
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	16777215,
	8388607,
	4194303,
	2097151,
	1048575,
	524287,
	262143,
	131071,
	65535,
	32767,
	16383,
	8191,
	4095,
	2047,
	1023,
	511,
	255,
	127,
	63,
	31,
	15,
	7,
	3,
	0,
	0,
}

func ipToInt32(ip net.IP) int32 {
	ip = ip.To4()
	return int32(ip[0])<<24 | int32(ip[1])<<16 | int32(ip[2])<<8 | int32(ip[3])
}

func int32ToIp(n int32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func GetBroadcast(network *net.IPNet) net.IP {
	if len(network.IP) != net.IPv4len {
		panic("utils: Not IPv4 network")
	}

	size, _ := network.Mask.Size()
	offset := offsets[size]
	if offset == 0 {
		panic("utils: Invalid network")
	}

	n := int(ipToInt32(network.IP[len(network.IP)-4:]))
	return int32ToIp(int32(n + offset))
}

func CopyIp(ip net.IP) (ipc net.IP) {
	ipc = make(net.IP, len(ip))
	copy(ipc, ip)
	return
}

func CopyNetwork(network *net.IPNet) (networkc *net.IPNet) {
	networkc = &net.IPNet{
		IP:   make(net.IP, len(network.IP)),
		Mask: make(net.IPMask, len(network.Mask)),
	}
	copy(networkc.IP, network.IP)
	copy(networkc.Mask, network.Mask)
	return
}

func IncIp(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func IterNetwork(network *net.IPNet) <-chan net.IP {
	iter := make(chan net.IP)

	go func() {
		ip := network.IP
		prev := make(net.IP, len(ip))
		i := 0
		for ip := ip.Mask(network.Mask); network.Contains(ip); IncIp(ip) {
			if prev != nil {
				if i < 2 {
					i += 1
				} else {
					iter <- prev
					prev = make(net.IP, len(ip))
				}
			}
			copy(prev, ip)
		}
		close(iter)
	}()

	return iter
}

func GetGateways() (gateways map[string]net.IP, err error) {
	gateways = map[string]net.IP{}
	gatewaysList := map[string]*list.List{}
	gatewaySets := map[string]set.Set{}
	nilAddr := []byte{0, 0, 0, 0}

	output, err := ExecOutput("", "route", "-n")
	if err != nil {
		return
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 8 {
			addr := net.ParseIP(fields[1])
			if addr == nil || bytes.HasSuffix(addr, nilAddr) {
				continue
			}
			iface := fields[7]

			var gwList *list.List
			gwSet, ok := gatewaySets[iface]
			if !ok {
				gwSet = set.NewSet()
				gatewaySets[iface] = gwSet

				gwList = list.New()
				gatewaysList[iface] = gwList
			} else {
				if gwSet.Contains(addr) {
					continue
				}

				gwList = gatewaysList[iface]
			}

			if addr[len(addr)-1] == 1 {
				gwList.PushFront(addr)
			} else {
				gwList.PushBack(addr)
			}
		}
	}

	for iface, gwList := range gatewaysList {
		gateways[iface] = gwList.Front().Value.(net.IP)
	}

	return
}

func GetInterfaceAddr(iface string) (ifaceAddr *InterfaceAddr, err error) {
	gateways, err := GetGateways()
	if err != nil {
		return
	}

	ifaceAddr = &InterfaceAddr{}

	ifaces, err := net.Interfaces()
	if err != nil {
		err = constants.ReadError{
			errors.Wrap(err, "utils: Failed to read network interfaces"),
		}
		return
	}

	for _, itf := range ifaces {
		if itf.Name != iface {
			continue
		}

		addrs, e := itf.Addrs()
		if e != nil {
			err = constants.ReadError{
				errors.Wrap(e, "utils: Failed to read network addresses"),
			}
			return
		}

		for _, addr := range addrs {
			adr, network, e := net.ParseCIDR(addr.String())
			if e != nil {
				err = constants.UnknownError{
					errors.Wrap(e, "utils: Failed to parse network"),
				}
				return
			}

			if strings.Contains(adr.String(), ":") {
				if ifaceAddr.Network6 == nil {
					ifaceAddr.Address6 = adr
					ifaceAddr.Network6 = network
				}
			} else if ifaceAddr.Network == nil {
				ifaceAddr.Gateway = gateways[itf.Name]
				ifaceAddr.Address = adr
				ifaceAddr.Network = network
			}
		}
	}

	return
}
