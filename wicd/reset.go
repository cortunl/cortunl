package wicd

func Reset() (err error) {
	conn, err := daemon()
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
