package dbus

import (
	"github.com/dropbox/godropbox/errors"
	"github.com/guelfey/go.dbus"
)

type Call struct {
	call *dbus.Call
	Body []interface{}
}

func (c *Call) init() (err error) {
	if c.call.Err != nil {
		err = &CallError{
			errors.Wrap(err, "dbus: DBus call error"),
		}
		return
	}

	c.Body = c.call.Body

	return
}

func (c *Call) Bind(store ...interface{}) (err error) {
	err = c.call.Store(store...)
	if err != nil {
		err = &BindError{
			errors.Wrap(err, "dbus: DBus bind error"),
		}
		return
	}

	return
}
