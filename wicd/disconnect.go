package wicd

func Disconnect() (err error) {
	conn, err := daemon()
	if err != nil {
		return
	}

	_, err = conn.Call("Disconnect")
	if err != nil {
		return
	}

	return
}
