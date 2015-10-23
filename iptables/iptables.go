package iptables

import (
	"fmt"
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"net"
)

type IpTables struct {
	state    bool
	rules    [][]string
	rules6   [][]string
	Inputs   []*settings.Input
	Bridge   string
	Network  *net.IPNet
	Network6 *net.IPNet
}

func (t *IpTables) addRule(rule []string, ipv6 bool) {
	comment := []string{
		"-m", "comment",
		"--comment", fmt.Sprintf("cortunl_%s", utils.Uuid()),
	}
	rule = append(rule, comment...)

	if ipv6 {
		t.rules6 = append(t.rules6, rule)
	} else {
		t.rules = append(t.rules, rule)
	}
}

func (t *IpTables) Init() {
	t.rules = [][]string{}
	t.rules6 = [][]string{}

	for _, input := range t.Inputs {
		for i, network := range []*net.IPNet{t.Network, t.Network6} {
			ipv6 := i == 1

			if input.NatInterface {
				t.addRule([]string{
					"POSTROUTING",
					"-t", "nat",
					"-o", input.Interface,
					"-j", "MASQUERADE",
					"-s", network.String(),
				}, ipv6)
			}

			if input.NatInterface {
				t.addRule([]string{
					"FORWARD",
					"-i", input.Interface,
					"-o", t.Bridge,
					"-m", "state",
					"--state", "RELATED,ESTABLISHED",
					"-j", "ACCEPT",
				}, ipv6)
			} else {
				t.addRule([]string{
					"FORWARD",
					"-i", input.Interface,
					"-o", t.Bridge,
					"-j", "ACCEPT",
				}, ipv6)
			}

			if input.AllTraffic {
				t.addRule([]string{
					"FORWARD",
					"-i", t.Bridge,
					"-o", input.Interface,
					"-j", "ACCEPT",
				}, ipv6)
			} else {
				for _, network := range input.Networks {
					if ipv6 != utils.IsIPNet6(network) {
						continue
					}

					t.addRule([]string{
						"FORWARD",
						"-i", t.Bridge,
						"-o", input.Interface,
						"-d", network.String(),
						"-j", "ACCEPT",
					}, ipv6)
				}
			}
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

	err = utils.Exec("", "iptables", "-P", "FORWARD", "DROP")
	if err != nil {
		return
	}

	err = utils.Exec("", "ip6tables", "-P", "FORWARD", "DROP")
	if err != nil {
		return
	}

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
