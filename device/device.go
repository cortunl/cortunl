package device

const (
	Wired    DeviceType = "wired"
	Wireless DeviceType = "wireless"
)

type DeviceType string

type Device struct {
	Name string
	Type DeviceType
}
