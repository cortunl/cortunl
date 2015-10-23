package router

import (
	"github.com/dropbox/godropbox/errors"
)

type ConnectError struct {
	errors.DropboxError
}
