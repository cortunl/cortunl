package hostapd

const (
	Auto    = "auto"
	NetLink = "nl80211"
	Realtek = "rtl871xdrv"
)

const conf = `driver=%s
ssid=%s
interface=%s
channel=%d
auth_algs=1
wpa=2
wpa_passphrase=%s
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP`
