package dhcp

import (
	"bytes"
	"fmt"
	"github.com/cortunl/cortunl/constants"
	"github.com/cortunl/cortunl/utils"
	"github.com/dropbox/godropbox/errors"
	"net"
	"os/exec"
	"path/filepath"
)

type dhcp6 struct {
	cmd       *exec.Cmd
	path      string
	output    *bytes.Buffer
	Interface string
	Network6  *net.IPNet
}

func (d *dhcp6) writeConf() (err error) {
	d.path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	d.path = filepath.Join(d.path, confName6)

	data := fmt.Sprintf(conf,
		d.Interface,
		d.Network6.String(),
	)

	err = utils.CreateWrite(d.path, data)
	if err != nil {
		return
	}

	return
}

func (d *dhcp6) Start() (err error) {
	d.output = &bytes.Buffer{}

	d.cmd = exec.Command("radvd", "--config", d.path)
	d.cmd.Stdout = d.output
	d.cmd.Stderr = d.output

	err = d.cmd.Start()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Failed to exec radvd"),
		}
		return
	}

	return
}

func (d *dhcp6) Stop() (err error) {
	err = d.cmd.Process.Kill()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Failed to stop radvd"),
		}
		return
	}

	return
}

func (d *dhcp6) Wait() (err error) {
	err = d.cmd.Wait()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Radvd exec error"),
		}
		return
	}

	return
}
