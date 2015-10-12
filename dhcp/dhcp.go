package dhcp

import (
	"net"
)

type Dhcp struct {
	v4        *dhcp4
	v6        *dhcp6
	Interface string
	Network   *net.IPNet
	Network6  *net.IPNet
}

func (d *Dhcp) Init() {
	d.v4 = &dhcp4{
		Interface: d.Interface,
		Network:   d.Network,
	}

	d.v6 = &dhcp6{
		Interface: d.Interface,
		Network6:  d.Network6,
	}
}

func (d *Dhcp) Start() (err error) {
	d.v4.Start()
	if err != nil {
		return
	}

	d.v6.Start()
	if err != nil {
		return
	}

	return
}

func (d *Dhcp) Stop() (err error) {
	d.v4.Stop()
	if err != nil {
		return
	}

	d.v6.Stop()
	if err != nil {
		return
	}

	return
}

func (d *Dhcp) Wait() (err error) {
	d.v4.Wait()
	if err != nil {
		return
	}

	d.v6.Wait()
	if err != nil {
		return
	}

	return
}
