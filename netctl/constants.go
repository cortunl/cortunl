package netctl

const (
	confNamePrefix = "cortunl_"
	confPathPrefix = "/etc/netctl/" + confNamePrefix
)

const conf = `Description='Cortunl Connection'
Connection=%s
Interface=%s
IP=dhcp
`
