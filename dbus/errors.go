package dbus

import (
	"github.com/dropbox/godropbox/errors"
)

type ConnError struct {
	errors.DropboxError
}

type CallError struct {
	errors.DropboxError
}

type BindError struct {
	errors.DropboxError
}

type CloseError struct {
	errors.DropboxError
}
