package dbus

import (
	"fmt"
	"github.com/docker/docker/vendor/src/github.com/godbus/dbus"
	"github.com/dropbox/godropbox/errors"
)

type Conn struct {
	conn *dbus.Conn
	obj  *dbus.Object
}

func (c *Conn) Call(method string, args ...interface{}) (
	call *Call, err error) {

	call = &Call{
		call: c.obj.Call(method, 0, args...),
	}
	err = call.init()
	if err != nil {
		return
	}

	return
}

func (c *Conn) Close() (err error) {
	err = c.conn.Close()
	if err != nil {
		err = &CloseError{
			errors.Wrap(err, "dbus: DBus connection close error"),
		}
		return
	}

	return
}

func NewConn(typ BusType, dest string, path string) (conn *Conn, err error) {
	c := &dbus.Conn{}

	switch typ {
	case SessionBus:
		c, err = dbus.SessionBus()
	case SessionBusPrivate:
		c, err = dbus.SessionBusPrivate()
	case SystemBus:
		c, err = dbus.SystemBus()
	case SystemBusPrivate:
		c, err = dbus.SystemBusPrivate()
	default:
		panic(fmt.Sprintf("dbus: Unknown bus type %d", typ))
	}
	if err != nil {
		err = &ConnError{
			errors.Wrap(err, "dbus: DBus connection error"),
		}
		return
	}

	obj := c.Object(dest, dbus.ObjectPath(path))

	conn = &Conn{
		conn: c,
		obj:  obj,
	}

	return
}
