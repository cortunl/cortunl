package iptables

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
)

type IpTables struct {
	Input  string
	Output string
}

func (i *IpTables) getRules() (rules [][]string) {
	rules = [][]string{
		[]string{
			"-t", "nat",
			"-A", "POSTROUTING",
			"-o", i.Input,
			"-j", "MASQUERADE",
		},
		[]string{
			"-A", "FORWARD",
			"-i", i.Input,
			"-o", i.Output,
			"-m", "state",
			"--state", "RELATED,ESTABLISHED",
			"-j", "ACCEPT",
		},
		[]string{
			"-A", "FORWARD",
			"-i", i.Output,
			"-o", i.Input,
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
