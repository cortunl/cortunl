package wicd

func CancelConnect() (err error) {
	conn, err := daemon()
	if err != nil {
		return
	}

	_, err = conn.Call("CancelConnect")
	if err != nil {
		return
	}

	return
}
