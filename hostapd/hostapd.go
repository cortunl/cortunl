package hostapd

import (
	"bytes"
	"fmt"
	"github.com/cortunl/cortunl/constants"
	"github.com/cortunl/cortunl/utils"
	"github.com/dropbox/godropbox/errors"
	"os"
	"os/exec"
	"path/filepath"
)

type Hostapd struct {
	cmd       *exec.Cmd
	path      string
	output    *bytes.Buffer
	Driver    Driver
	Interface string
	Ssid      string
	Channel   int
	Password  string
}

func (h *Hostapd) writeConf() (err error) {
	h.path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	h.path = filepath.Join(h.path, confName)

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
	err = h.writeConf()
	if err != nil {
		return
	}

	h.output = &bytes.Buffer{}

	h.cmd = exec.Command("hostapd", h.path)
	h.cmd.Stdout = os.Stdout
	h.cmd.Stderr = os.Stdout

	err = h.cmd.Start()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "hostapd: Failed to exec"),
		}
		return
	}

	return
}

func (h *Hostapd) Stop() (err error) {
	if h.cmd == nil {
		return
	}

	err = h.cmd.Process.Kill()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "hostapd: Failed to stop exec"),
		}
		return
	}

	return
}

func (h *Hostapd) Wait() (err error) {
	if h.cmd == nil {
		return
	}

	err = h.cmd.Wait()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "hostapd: Exec error"),
		}
		return
	}

	return
}
