package iptables

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
)

type IpTables struct {
}

func (i *IpTables) AddRules() (err error) {
	err = utils.EnableIpForwarding()
	if err != nil {
		return
	}

	rules := [][]string{
		[]string{
			"-t", "nat",
			"-A", "POSTROUTING",
			"-o", "eth0",
			"-j", "MASQUERADE",
		},
		[]string{
			"-A", "FORWARD",
			"-i", "eth0",
			"-o", "wlan0",
			"-m", "state",
			"--state", "RELATED,ESTABLISHED",
			"-j", "ACCEPT",
		},
		[]string{
			"-A", "FORWARD",
			"-i", "wlan0",
			"-o", "eth0",
			"-j", "ACCEPT",
		},
	}

	for _, rule := range rules {
		comment := []string{
			"-m", "comment",
			"--comment", fmt.Sprintf("cortunl_%s", utils.Uuid()),
		}

		rule = append(rule, comment...)
	}

	return
}
