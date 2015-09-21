package dhcpd

const (
	confName = "dhcpd.conf"
)

const conf = `subnet %s netmask %s {
	range %s %s;
	option subnet-mask %s;
	option broadcast-address %s;
	option routers %s;
	option domain-name "%s";
	option domain-name-servers %s;
	default-lease-time 172800;
	max-lease-time 172800;
}`
