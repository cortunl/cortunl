package hostapd

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
	"path/filepath"
)

type Hostapd struct {
	path      string
	Driver    Driver
	Ssid      string
	Interface string
	Channel   int
	Password  string
}

func (h *Hostapd) writeConf() (err error) {
	h.path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	h.path = filepath.Join(h.path, "hostapd.conf")

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

	data := fmt.Sprintf(conf, driver, h.Ssid, h.Interface, h.Channel, wpaData)

	err = utils.CreateWrite(h.path, data)
	if err != nil {
		return
	}

	return
}

func (h *Hostapd) Start() (err error) {
	err = utils.Exec("hostapd", h.path)
	if err != nil {
		return
	}

	return
}
