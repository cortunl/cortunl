package wicd

func Reset() (err error) {
	conn, err := daemon()
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Call("SetSignalDisplayType", "0")
	if err != nil {
		return
	}

	return
}
