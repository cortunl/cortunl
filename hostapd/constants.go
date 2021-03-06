package hostapd

type Driver string

const (
	AutoChan int    = 0
	AutoDrv  Driver = "auto"
	NetLink  Driver = "netlink"
	Realtek  Driver = "realtek"
	confName        = "hostapd.conf"
)

const conf = `driver=%s
ssid=%s
interface=%s
bridge=%s
hw_mode=g
country_code=US
ieee80211d=1
wmm_enabled=1
ieee80211n=1
#ieee80211ac=1
channel=%d%s`

const wpaConf = `
auth_algs=3
wpa=3
wpa_passphrase=%s
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP`
