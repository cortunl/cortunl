package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cortunl/cortunl/netctl"
	"github.com/cortunl/cortunl/network"
	"github.com/cortunl/cortunl/router"
	"github.com/cortunl/cortunl/settings"
	"os"
	"os/signal"
	"syscall"
)

var (
	scanWirelessCmd = cli.Command{
		Name:   "scan",
		Action: scanWireless,
	}
	connectWiredCmd = cli.Command{
		Name:   "connect-wired",
		Action: connectWired,
	}
	connectWirelessCmd = cli.Command{
		Name:   "connect-wireless",
		Action: connectWireless,
	}
	disconnectCmd = cli.Command{
		Name:   "disconnect",
		Action: disconnect,
	}
	statusCmd = cli.Command{
		Name:   "status",
		Action: status,
	}
	routerCmd = cli.Command{
		Name:   "router",
		Action: routr,
	}
)

func init() {
	commands = append(commands, scanWirelessCmd)
	commands = append(commands, connectWiredCmd)
	commands = append(commands, connectWirelessCmd)
	commands = append(commands, disconnectCmd)
	commands = append(commands, statusCmd)
	commands = append(commands, routerCmd)
}

func scanWireless(c *cli.Context) {
	networks, err := netctl.GetNetworks("wlan0")
	if err != nil {
		panic(err)
	}

	for _, net := range networks {
		fmt.Printf("%s: %s\n", net.Ssid, net.Security.Type())
	}
}

func connectWired(c *cli.Context) {
	net := &network.WiredNetwork{
		&network.Network{
			Interface: "eth0",
		},
	}

	err := netctl.Connect(net)
	if err != nil {
		panic(err)
	}
}

func connectWireless(c *cli.Context) {
	networks, err := netctl.GetNetworks(
		settings.Settings.Routers[0].Inputs[0].Interface)
	if err != nil {
		panic(err)
	}

	for _, net := range networks {
		if net.Ssid == "arch" {
			err = net.Security.Set("password", "")
			if err != nil {
				panic(err)
			}

			err = netctl.Connect(net)
			if err != nil {
				panic(err)
			}
		}
	}
}

func disconnect(c *cli.Context) {
	err := netctl.DisconnectAll()
	if err != nil {
		panic(err)
	}
}

func status(c *cli.Context) {
	status, err := netctl.Status("eth0")
	if err != nil {
		panic(err)
	}

	fmt.Printf("State: %s\n", status)
}

func routr(c *cli.Context) {
	rter := &router.Router{
		Settings: settings.Settings.Routers[0],
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		err := rter.Stop()
		if err != nil {
			panic(err)
		}
	}()

	err := rter.Run()
	if err != nil {
		e := rter.Stop()
		if e != nil {
			panic(e)
		}
		panic(err)
	}
}
