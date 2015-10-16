package netctl

import (
	"github.com/dropbox/godropbox/errors"
)

type ParseError struct {
	errors.DropboxError
}

type NotFound struct {
	errors.DropboxError
}
