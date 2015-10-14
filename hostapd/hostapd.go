package hostapd

import (
	"fmt"
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

func (h *Hostapd) writeConf() (path string, err error) {
	path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	path = filepath.Join(path, confName)

	driver := ""
	switch h.Driver {
	case NetLink:
		driver = "nl80211"
	case Realtek:
		driver = "rtl871xdrv"
	default:
		driver = "nl80211"
	}

	wpaData := ""
	if h.Password != "" {
		wpaData = fmt.Sprintf(wpaConf, h.Password)
	}

	data := fmt.Sprintf(conf,
		driver,
		h.Ssid,
		h.Interface,
		h.Bridge,
		h.Channel,
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
