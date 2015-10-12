package dhcp

import (
	"bytes"
	"fmt"
	"github.com/cortunl/cortunl/constants"
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"github.com/dropbox/godropbox/errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type dhcp4 struct {
	cmd       *exec.Cmd
	path      string
	output    *bytes.Buffer
	Interface string
	Network   *net.IPNet
}

func (d *dhcp4) writeConf() (err error) {
	d.path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	d.path = filepath.Join(d.path, confName)

	broadcast := utils.GetBroadcast(d.Network).String()
	mask := net.IP(d.Network.Mask).String()

	router := ""
	start := ""
	end := ""
	for ip := range utils.IterNetwork(d.Network) {
		if router == "" {
			router = ip.String()
		} else if start == "" {
			start = ip.String()
		} else {
			end = ip.String()
		}
	}

	data := fmt.Sprintf(conf,
		d.Network.IP.String(),
		mask,
		start,
		end,
		mask,
		broadcast,
		router,
		strings.Join(settings.Settings.DnsServers, ", "),
	)

	err = utils.CreateWrite(d.path, data)
	if err != nil {
		return
	}

	return
}

func (d *dhcp4) Start() (err error) {
	err = d.writeConf()
	if err != nil {
		return
	}

	d.output = &bytes.Buffer{}

	d.cmd = exec.Command("dhcpcd", "--config", d.path)
	d.cmd.Stdout = d.output
	d.cmd.Stderr = d.output

	err = d.cmd.Start()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Failed to exec dhcpcd"),
		}
		return
	}

	return
}

func (d *dhcp4) Stop() (err error) {
	if d.cmd == nil {
		return
	}

	err = d.cmd.Process.Kill()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Failed to stop dhcpcd"),
		}
		return
	}

	err = d.Wait()
	if err != nil {
		return
	}

	return
}

func (d *dhcp4) Wait() (err error) {
	if d.cmd == nil {
		return
	}

	err = d.cmd.Wait()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Dhcpcd exec error"),
		}
		return
	}

	return
}
