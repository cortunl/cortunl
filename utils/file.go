package utils

import (
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
	"os"
)

func Create(path string) (file *os.File, err error) {
	file, err = os.Create(path)
	if err != nil {
		err = &constants.WriteError{
			errors.Wrapf(err, "utils: Failed to create '%s'", path),
		}
		return
	}

	return
}

func CreateWrite(path string, data string) (err error) {
	file, err := Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		err = &constants.WriteError{
			errors.Wrapf(err, "utils: Failed to write to file '%s'", path),
		}
		return
	}

	return
}
