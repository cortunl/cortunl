package wicd

func Disconnect() (err error) {
	conn, err := daemon()
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = conn.Call("Disconnect")
	if err != nil {
		return
	}

	return
}
