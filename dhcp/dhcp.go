package dhcp

import (
	"fmt"
	"github.com/cortunl/cortunl/runner"
	"github.com/cortunl/cortunl/utils"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Dhcp struct {
	runner.Runner
	Bridge      string
	LocalDomain string
	DnsServers  []string
	DnsServers6 []string
	Network     *net.IPNet
	Network6    *net.IPNet
}

func (d *Dhcp) writeConf() (path string, err error) {
	path, err = utils.GetTempDir()
	if err != nil {
		return
	}
	path = filepath.Join(path, confName)

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
	dnsServers = append(dnsServers, d.DnsServers...)
	dnsServers = append(dnsServers, d.DnsServers6...)
	for _, svr := range dnsServers {
		servers += fmt.Sprintf("server=%s\n", svr)
	}

	data := fmt.Sprintf(conf,
		servers,
		d.Bridge,
		d.LocalDomain,
		d.LocalDomain,
		start,
		end,
		mask,
		network6,
		cidr6,
		router,
		router,
		d.DnsServers[0], // TODO
		router6,
		d.DnsServers6[0], // TODO
		d.LocalDomain,
	)

	err = utils.Write(path, data)
	if err != nil {
		return
	}

	return
}

func (d *Dhcp) Start() (err error) {
	path, err := d.writeConf()
	if err != nil {
		return
	}

	cmd := exec.Command("dnsmasq", "--keep-in-foreground",
		fmt.Sprintf("--conf-file=%s", path))
	err = d.Run(cmd, func() {
		os.RemoveAll(filepath.Dir(path))
	})
	if err != nil {
		return
	}

	return
}
