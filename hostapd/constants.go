package hostapd

type Driver string

const (
	Auto     Driver = "auto"
	NetLink  Driver = "netlink"
	Realtek  Driver = "realtek"
	confName        = "hostapd.conf"
)

const conf = `driver=%s
ssid=%s
interface=%s
channel=%d%s`

const wpaConf = `
auth_algs=1
wpa=2
wpa_passphrase=%s
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP`
