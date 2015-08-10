package main

import (
	"github.com/cortunl/cortunl/cmd"
)

func main() {
	_, err := cmd.Run()
	if err != nil {
		return
	}
}
