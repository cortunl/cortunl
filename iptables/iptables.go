package iptables

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
	"net"
)

type IpTables struct {
	state    bool
	rules    [][]string
	rules6   [][]string
	Input    string
	Output   string
	Network  *net.IPNet
	Network6 *net.IPNet
}

func (t *IpTables) Init() {
	for i, network := range []*net.IPNet{t.Network, t.Network6} {
		rules := [][]string{
			[]string{
				"POSTROUTING",
				"-t", "nat",
				"-o", t.Input,
				"-j", "MASQUERADE",
				"-s", network.String(),
			},
			[]string{
				"FORWARD",
				"-i", t.Input,
				"-o", t.Output,
				"-m", "state",
				"--state", "RELATED,ESTABLISHED",
				"-j", "ACCEPT",
			},
			[]string{
				"FORWARD",
				"-i", t.Output,
				"-o", t.Input,
				"-j", "ACCEPT",
			},
		}

		for i, rule := range rules {
			comment := []string{
				"-m", "comment",
				"--comment", fmt.Sprintf("cortunl_%s", utils.Uuid()),
			}

			rules[i] = append(rule, comment...)
		}

		if i == 0 {
			t.rules = rules
		} else {
			t.rules6 = rules
		}
	}

	return
}

func (t *IpTables) run(mode string, force bool) (err error) {
	for i, exec := range []string{"iptables", "ip6tables"} {
		var rules [][]string

		if i == 0 {
			rules = t.rules
		} else {
			rules = t.rules6
		}

		for _, rule := range rules {
			args := append([]string{mode}, rule...)

			err = utils.Exec("", exec, args...)
			if err != nil {
				if force {
					err = nil
				} else {
					return
				}
			}
		}
	}

	return
}

func (t *IpTables) AddRules() (err error) {
	if t.state {
		return
	}
	t.state = true

	err = t.run("-I", false)
	if err != nil {
		return
	}

	return
}

func (t *IpTables) RemoveRules() (err error) {
	if !t.state {
		return
	}
	t.state = false

	err = t.run("-D", true)
	if err != nil {
		return
	}

	return
}
