package cmd

import (
	"github.com/dropbox/godropbox/errors"
)

type RunError struct {
	errors.DropboxError
}
