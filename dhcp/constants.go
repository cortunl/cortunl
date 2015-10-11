package dhcp

const (
	confName  = "dhcpd.conf"
	confName6 = "radvd.conf"
)

const conf = `subnet %s netmask %s {
	range %s %s;
	option subnet-mask %s;
	option broadcast-address %s;
	option routers %s;
	option domain-name-servers %s;
	default-lease-time 172800;
	max-lease-time 172800;
}`

const conf6 = `interface %s {
	AdvSendAdvert on;
	MinRtrAdvInterval 3;
	MaxRtrAdvInterval 10;
	prefix %s {
		AdvOnLink on;
		AdvAutonomous on;
		AdvRouterAddr on;
	};
};`
