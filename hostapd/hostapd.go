package hostapd

import (
	"fmt"
	"github.com/cortunl/cortunl/utils"
)

type Hostapd struct {
	Path      string
	Driver    Driver
	Ssid      string
	Interface string
	Channel   int
	Password  string
}

func (h *Hostapd) writeConf() (err error) {
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

	err = utils.CreateWrite(h.Path, data)
	if err != nil {
		return
	}

	return
}
