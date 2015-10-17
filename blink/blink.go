package blink

import (
	"github.com/cortunl/cortunl/settings"
	"github.com/cortunl/cortunl/utils"
	"github.com/dropbox/godropbox/container/set"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	supported   set.Set
	unsupported set.Set
)

func init() {
	supported = set.NewSet()
	unsupported = set.NewSet()
}

func Init() {
	waiter := sync.WaitGroup{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range ifaces {
		waiter.Add(1)
		go func(iface string) {
			check(iface)
			waiter.Done()
		}(iface.Name)
	}

	waiter.Wait()
}

func check(iface string) bool {
	err := utils.ExecSilent("", "ethtool", "-p", iface, "1")
	if err != nil {
		unsupported.Add(iface)
		return false
	}
	supported.Add(iface)
	return true
}

func CanBlink(iface string) bool {
	if supported.Contains(iface) {
		return true
	} else if unsupported.Contains(iface) {
		return false
	}
	return check(iface)
}

func Blink(iface string) (err error) {
	go func() {
		err = utils.Exec("", "ethtool", "-p", iface, strconv.Itoa(
			int(settings.Settings.BlinkDuration/time.Second)))
	}()

	time.Sleep(50 * time.Millisecond)
	return
}
