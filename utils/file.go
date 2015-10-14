package utils

import (
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Read(path string) (data string, err error) {
	dataByt, err := ioutil.ReadFile(path)
	if err != nil {
		err = &constants.ReadError{
			errors.Wrapf(err, "utils: Failed to read '%s'", path),
		}
		return
	}
	data = string(dataByt)

	return
}

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

func MkdirAll(path string) (err error) {
	err = os.MkdirAll(path, 0755)
	if err != nil {
		err = &constants.WriteError{
			errors.Wrapf(err, "utils: Failed to create '%s'", path),
		}
		return
	}

	return
}

func Write(path string, data string) (err error) {
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

func GetTempDir() (path string, err error) {
	path = filepath.Join(
		constants.TempDir,
		string(filepath.Separator),
		RandStr(16),
	)

	err = MkdirAll(path)
	if err != nil {
		return
	}

	return
}
