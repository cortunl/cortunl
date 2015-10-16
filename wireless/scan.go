package wireless

import (
	"github.com/cortunl/cortunl/network"
	"github.com/cortunl/cortunl/security"
	"github.com/cortunl/cortunl/utils"
	"github.com/dropbox/godropbox/errors"
	"math"
	"strconv"
	"strings"
)

func parseField(line, field string) (val string, err error) {
	split := strings.Split(line, field)
	if len(split) != 2 {
		err = &ParseError{
			errors.Newf("wireless: Failed to parse wireless field '%s'",
				field[:len(field)-1]),
		}
		return
	}

	val = strings.TrimSpace(split[1])
	return
}

func GetNetworks(iface string) (networks []*network.WirelessNetwork,
	err error) {

	err = utils.Exec("", "ip", "link", "set", iface, "up")
	if err != nil {
		return
	}

	networks = []*network.WirelessNetwork{}

	output, err := utils.ExecOutput("", "iwlist", iface, "scan")
	if err != nil {
		return
	}

	var netwk *network.WirelessNetwork
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Cell ") {
			val, e := parseField(line, "Address:")
			if e != nil {
				err = e
				return
			}

			netwk = &network.WirelessNetwork{
				Mac: val,
			}
			networks = append(networks, netwk)
		} else if strings.HasPrefix(line, "Channel:") {
			val, e := parseField(line, "Channel:")
			if e != nil {
				err = e
				return
			}

			valN, e := strconv.Atoi(val)
			if e != nil {
				err = &ParseError{
					errors.Newf("wireless: Failed to parse "+
						"wireless channel '%s'", val),
				}
				return
			}

			netwk.Channel = valN
		} else if strings.HasPrefix(line, "Frequency:") {
			if !strings.Contains(line, "Channel") {
				continue
			}

			fields := strings.Split(line, "Channel ")
			if len(fields) != 2 {
				continue
			}

			val := fields[1]
			val = val[:len(val)-1]

			valN, e := strconv.Atoi(val)
			if e != nil {
				continue
			}

			netwk.Channel = valN
		} else if strings.HasPrefix(line, "Quality=") {
			val, e := parseField(line, "Quality=")
			if e != nil {
				err = e
				return
			}

			val = strings.Fields(val)[0]
			split := strings.Split(val, "/")
			if len(split) != 2 {
				err = &ParseError{
					errors.Newf("wireless: Failed to parse "+
						"wireless quality '%s'", val),
				}
				return
			}

			num, e := strconv.Atoi(split[0])
			if e != nil {
				err = &ParseError{
					errors.Newf("wireless: Failed to parse "+
						"wireless quality '%s'", val),
				}
				return
			}

			den, e := strconv.Atoi(split[1])
			if e != nil {
				err = &ParseError{
					errors.Newf("wireless: Failed to parse "+
						"wireless quality '%s'", val),
				}
				return
			}

			quality := float64(num) / float64(den)

			netwk.Quality = int(math.Floor(quality * 100))
		} else if strings.HasPrefix(line, "ESSID:") {
			val, e := parseField(line, "ESSID:")
			if e != nil {
				err = e
				return
			}

			netwk.Ssid = val[1 : len(val)-1]
		} else if line == "Encryption key:off" {
			netwk.Security = security.GetSecurity("open")
		} else if strings.Contains(line, "IE") &&
			strings.Contains(line, "WPA") {

			netwk.Security = security.GetSecurity("wpa")
		}
	}

	return
}
