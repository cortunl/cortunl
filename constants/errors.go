package constants

import (
	"github.com/dropbox/godropbox/errors"
)

type ReadError struct {
	errors.DropboxError
}

type WriteError struct {
	errors.DropboxError
}

type UnknownError struct {
	errors.DropboxError
}

type ExecError struct {
	errors.DropboxError
}
