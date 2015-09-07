package wicd

import (
	"github.com/Sirupsen/logrus"
	"github.com/cortunl/cortunl/network"
	"github.com/cortunl/cortunl/security"
	"github.com/cortunl/cortunl/utils"
	"github.com/dropbox/godropbox/errors"
	"strconv"
	"strings"
	"sync"
)

var (
	lock = sync.Mutex{}
)

func getNetworkNum(ssid string) (num string, err error) {
	output, err := utils.ExecOutput("", "wicd-cli", "--wireless",
		"--scan", "--list-networks")
	if err != nil {
		return
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			err = &ParseError{
				errors.Newf("wicd: Too few fields on line '%s'", line),
			}
			return
		}

		if fields[3] == ssid {
			num = fields[0]
			return
		}
	}

	err = &NotFound{
		errors.Newf("wicd: Failed to find ssid '%s'", ssid),
	}
	return
}

func GetNetworks() (networks []*network.WirelessNetwork, err error) {
	lock.Lock()
	defer lock.Unlock()

	output, err := utils.ExecOutput("", "wicd-cli", "--wireless",
		"--scan", "--list-networks")
	if err != nil {
		return
	}

	lines := strings.Split(output, "\n")
Network:
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}

		num := strings.Fields(line)[0]

		output, e := utils.ExecOutput("", "wicd-cli", "--wireless",
			"--network", num, "--network-details")
		if e != nil {
			err = e
			return
		}

		net := &network.WirelessNetwork{
			Network: &network.Network{
				Type: network.Wireless,
			},
		}

		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}

			fields := strings.Split(line, ":")
			if len(fields) < 2 {
				logrus.WithFields(logrus.Fields{
					"id":   num,
					"line": line,
				}).Error("wicd: Failed to parse line")
				continue Network
			}

			key := strings.TrimSpace(fields[0])
			val := strings.TrimSpace(fields[1])

			switch key {
			case "Essid":
				net.Ssid = val
			case "Encryption Method":
				sec := security.GetSecurity(strings.ToLower(val))
				if sec == nil {
					logrus.WithFields(logrus.Fields{
						"id":       num,
						"ssid":     net.Ssid,
						"security": val,
					}).Warning("wicd: Unknown encryption type")
					continue Network
				}

				net.Security = sec
			case "Quality":
				net.Quality, err = strconv.Atoi(val)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"id":      num,
						"ssid":    net.Ssid,
						"quality": val,
						"error":   err,
					}).Warning("wicd: Failed to parse quality")
					err = nil
					continue Network
				}
			case "Channel":
				net.Channel, err = strconv.Atoi(val)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"id":      num,
						"ssid":    net.Ssid,
						"channel": val,
						"error":   err,
					}).Warning("wicd: Failed to parse channel")
					err = nil
					continue Network
				}
			}
		}

		networks = append(networks, net)
	}

	return
}
