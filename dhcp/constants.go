package dhcp

const (
	confName  = "dnsmasq.conf"
)

const conf = `no-resolv
domain-needed
%sinterface=%s
no-hosts
dhcp-range=%s,%s,%s,12h
dhcp-range=%s,ra-stateless,ra-names,%s
dhcp-option=option:router,%s
dhcp-option=option:dns-server,%s,%s
dhcp-option=option6:dns-server,[%s],[%s]
dhcp-authoritative`
