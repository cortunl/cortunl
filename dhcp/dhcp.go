package dhcp

import (
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

type Dhcp struct {
	cmd       *exec.Cmd
	path      string
	Interface string
	Network   *net.IPNet
	Network6  *net.IPNet
}

func (d *Dhcp) writeConf() (err error) {
	d.path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	d.path = filepath.Join(d.path, confName)

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

	router6 := utils.CopyIp(d.Network6.IP)
	utils.IncIp(router6)
	networkSplit6 := strings.Split(d.Network6.String(), "/")
	network6 := networkSplit6[0]
	cidr6 := networkSplit6[1]

	servers := ""
	dnsServers := []string{}
	dnsServers = append(dnsServers, settings.Settings.DnsServers...)
	dnsServers = append(dnsServers, settings.Settings.DnsServers6...)
	for _, svr := range dnsServers {
		servers += fmt.Sprintf("server=%s\n", svr)
	}

	data := fmt.Sprintf(conf,
		servers,
		d.Interface,
		start,
		end,
		mask,
		network6,
		cidr6,
		router,
		router,
		settings.Settings.DnsServers[0],
		router6,
		settings.Settings.DnsServers6[0],
	)

	err = utils.CreateWrite(d.path, data)
	if err != nil {
		return
	}

	return
}

func (d *Dhcp) Start() (err error) {
	err = d.writeConf()
	if err != nil {
		return
	}

	d.cmd = exec.Command("dnsmasq", "--keep-in-foreground",
		fmt.Sprintf("--conf-file=%s", d.path))
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stdout

	err = d.cmd.Start()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Failed to exec dnsmasq"),
		}
		return
	}

	return
}

func (d *Dhcp) Stop() (err error) {
	if d.cmd == nil {
		return
	}

	err = d.cmd.Process.Kill()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Failed to stop dnsmasq"),
		}
		return
	}

	return
}

func (d *Dhcp) Wait() (err error) {
	if d.cmd == nil {
		return
	}

	err = d.cmd.Wait()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "dhcp: Dnsmasq exec error"),
		}
		return
	}

	return
}
