package wicd

import (
	"github.com/cortunl/cortunl/settings"
)

func Reset() (err error) {
	conn, err := daemon()
	if err != nil {
		return
	}

	wired := ""
	if len(settings.Settings.ExternalWired) > 0 {
		wired = settings.Settings.ExternalWired[0]
	}
	_, err = conn.Call("SetWiredInterface", wired)
	if err != nil {
		return
	}

	wireless := ""
	if len(settings.Settings.ExternalWireless) > 0 {
		wireless = settings.Settings.ExternalWireless[0]
	}
	_, err = conn.Call("SetWirelessInterface", wireless)
	if err != nil {
		return
	}

	_, err = conn.Call("SetForcedDisconnect", true)
	if err != nil {
		return
	}

	_, err = conn.Call("SetUseGlobalDNS", false)
	if err != nil {
		return
	}

	_, err = conn.Call("SetSignalDisplayType", uint32(0))
	if err != nil {
		return
	}

	_, err = conn.Call("SetAutoReconnect", false)
	if err != nil {
		return
	}

	_, err = conn.Call("SetWiredAutoConnectMethod", uint32(1))
	if err != nil {
		return
	}

	_, err = conn.Call("SetPreferWiredNetwork", false)
	if err != nil {
		return
	}

	return
}
