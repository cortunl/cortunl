package utils

import (
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
	"os"
	"os/exec"
)

func Exec(dir, name string, arg ...string) (err error) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if dir != "" {
		cmd.Dir = dir
	}

	err = cmd.Run()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrapf(err, "utils: Failed to exec '%s'", name),
		}
		return
	}

	return
}

func ExecOutput(dir, name string, arg ...string) (output string, err error) {
	cmd := exec.Command(name, arg...)
	if dir != "" {
		cmd.Dir = dir
	}

	outputByt, err := cmd.CombinedOutput()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrapf(err, "utils: Failed to exec '%s'", name),
		}
		return
	}
	output = string(outputByt)

	return
}
