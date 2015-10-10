package dhcp

type Dhcp interface {
	Start() error
	Stop() error
	Wait() error
}
