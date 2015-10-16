package hostapd

import (
	"fmt"
	"github.com/cortunl/cortunl/runner"
	"github.com/cortunl/cortunl/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		if runtime.GOARCH == "arm" {
			drv = Realtek
		} else {
			drv = NetLink
		}
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
		channel = 6 // TODO
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

	bin := ""
	switch h.getDriver() {
	case NetLink:
		bin = "hostapd"
	case Realtek:
		bin = "hostapd_rtl"
	}

	cmd := exec.Command(bin, path)
	err = h.Run(cmd, func() {
		os.RemoveAll(filepath.Dir(path))
	})
	if err != nil {
		return
	}

	return
}
