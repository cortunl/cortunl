package iptables

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
)

type IpTables struct {
	Input    string
	Output   string
	Network  string
	Network6 string
	rules    [][]string
	rules6   [][]string
}

func (t *IpTables) Init() {
	for i, network := range []string{t.Network, t.Network6} {
		rules := [][]string{
			[]string{
				"POSTROUTING",
				"-t", "nat",
				"-o", t.Input,
				"-j", "MASQUERADE",
				"-s", network,
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

func (t *IpTables) run(mode string) (err error) {
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
				return
			}
		}
	}

	return
}

func (t *IpTables) AddRules() (err error) {
	err = t.run("-I")
	if err != nil {
		return
	}

	return
}

func (t *IpTables) RemoveRules() (err error) {
	err = t.run("-D")
	if err != nil {
		return
	}

	return
}
