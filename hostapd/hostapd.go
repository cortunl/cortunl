package hostapd

import (
	"fmt"
	"github.com/cortunl/cortunl/netctl"
	"github.com/cortunl/cortunl/runner"
	"github.com/cortunl/cortunl/utils"
	"os"
	"os/exec"
	"path/filepath"
)

type Hostapd struct {
	runner.Runner
	Driver    Driver
	Bridge    string
	Interface string
	Ssid      string
	Channel   int
	Password  string
}

func (h *Hostapd) getDriver() (drv Driver) {
	switch h.Driver {
	case NetLink:
		drv = NetLink
	case Realtek:
		drv = Realtek
	default:
		drv = NetLink
	}

	return
}

func (h *Hostapd) writeConf() (path string, err error) {
	path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	path = filepath.Join(path, confName)

	driver := ""
	switch h.getDriver() {
	case NetLink:
		driver = "nl80211"
	case Realtek:
		driver = "rtl871xdrv"
	}

	wpaData := ""
	if h.Password != "" {
		wpaData = fmt.Sprintf(wpaConf, h.Password)
	}

	channel := h.Channel
	if channel == AutoChan {
		channels := map[int]int{}

		networks, e := netctl.GetNetworks(h.Interface)
		if e != nil {
			err = e
			return
		}

		for _, net := range networks {
			channels[net.Channel] += net.Quality
		}

		bestStrength := 0
		bestChannel := 0
		for channel := 1; channel < 12; channel++ {
			strength := 0

			strength += channels[channel-2]
			strength += channels[channel-1]
			strength += channels[channel]
			strength += channels[channel+1]
			strength += channels[channel+2]

			if channel != 1 && channel != 6 && channel != 11 {
				strength += 10
			}

			if bestChannel == 0 || strength <= bestStrength {
				bestStrength = strength
				bestChannel = channel
			}
		}

		channel = bestChannel
	}

	data := fmt.Sprintf(conf,
		driver,
		h.Ssid,
		h.Interface,
		h.Bridge,
		channel,
		wpaData,
	)

	err = utils.Write(path, data)
	if err != nil {
		return
	}

	return
}

func (h *Hostapd) Start() (err error) {
	path, err := h.writeConf()
	if err != nil {
		return
	}

	cmd := exec.Command("hostapd", path)
	err = h.Run(cmd, func() {
		os.RemoveAll(filepath.Dir(path))
	})
	if err != nil {
		return
	}

	return
}
