package iptables

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
)

type IpTables struct {
	Input   string
	Output  string
	Network string
	rules   [][]string
}

func (t *IpTables) generateRules() {
	t.rules = [][]string{
		[]string{
			"POSTROUTING",
			"-t", "nat",
			"-o", t.Input,
			"-j", "MASQUERADE",
			"-s", t.Network,
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

	for i, rule := range t.rules {
		comment := []string{
			"-m", "comment",
			"--comment", fmt.Sprintf("cortunl_%s", utils.Uuid()),
		}

		t.rules[i] = append(rule, comment...)
	}

	return
}

func (t *IpTables) AddRules() (err error) {
	t.generateRules()

	for _, rule := range t.rules {
		args := append([]string{"-I"}, rule...)

		err = utils.Exec("", "iptables", args...)
		if err != nil {
			return
		}
	}

	return
}
