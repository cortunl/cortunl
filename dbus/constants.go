package dbus

const (
	SessionBus        BusType = 0
	SessionBusPrivate BusType = 1
	SystemBus         BusType = 2
	SystemBusPrivate  BusType = 3
)

type BusType int
