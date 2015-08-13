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
				err = &ParseError{
					errors.Newf("wicd: Failed to parse line '%s'", line),
				}
				return
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
					}).Info("wicd: Unknown encryption type")
					continue Network
				}
			case "Quality:":
				net.Quality, err = strconv.Atoi(val)
				if err != nil {
					err = &ParseError{
						errors.Newf("wicd: Failed to parse quality '%s'",
							val),
					}
					return
				}
			case "Channel:":
				net.Channel, err = strconv.Atoi(val)
				if err != nil {
					err = &ParseError{
						errors.Newf("wicd: Failed to parse channel '%s'",
							val),
					}
					return
				}
			}
		}

		networks = append(networks, net)
	}

	return
}
