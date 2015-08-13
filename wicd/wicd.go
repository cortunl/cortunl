package wicd

import (
	"github.com/Sirupsen/logrus"
	"github.com/cortunl/cortunl/network"
	"github.com/dropbox/godropbox/errors"
	"github.com/pacur/pacur/utils"
	"strconv"
	"strings"
)

func Scan() (networks []*network.WirelessNetwork, err error) {
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
			"--network-details", "--network", num)
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

			fields := strings.Fields(line)
			if len(fields) < 2 {
				logrus.WithFields(logrus.Fields{
					"id":   num,
					"line": line,
				}).Error("wicd: Failed to parse line")
				continue Network
			}

			key := fields[0]
			val := fields[1]

			switch key {
			case "Essid:":
				net.Ssid = val
			case "Encryption Method:":
				switch val {
				case "WEP":
					net.Security = network.Wep
				case "WPA":
					net.Security = network.Wpa
				case "WPA2":
					net.Security = network.Wpa2
				case "NONE":
					net.Security = network.None
				default:
					logrus.WithFields(logrus.Fields{
						"id":       num,
						"ssid":     net.Ssid,
						"security": val,
					}).Warning("wicd: Unknown encryption type")
					continue Network
				}
			case "Quality:":
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
			case "Channel:":
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
