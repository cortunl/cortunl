package security

import (
	"github.com/dropbox/godropbox/errors"
)

type UnknownError struct {
	errors.DropboxError
}

type InvalidError struct {
	errors.DropboxError
}
