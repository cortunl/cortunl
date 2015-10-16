package netctl

const (
	confName = "cortunl"
	confPath = "/etc/netctl/" + confName
)

const conf = `Description='Cortunl Connection'
Connection=%s
Interface=%s
IP=dhcp
`
