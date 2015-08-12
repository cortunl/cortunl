package utils

import (
	"github.com/dropbox/godropbox/errors"
)

type ExecError struct {
	errors.DropboxError
}
