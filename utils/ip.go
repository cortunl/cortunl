package utils

import (
	"net"
)

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