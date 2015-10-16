package security

type Security interface {
	Init()
	Type() string
	Params() []*SecurityParam
	Set(string, string) error
	Properties() map[string]string
}

func GetSecurity(typ string) (intf Security) {
	switch typ {
	case "open":
		intf = &WpaSecurity{}
	case "wpa", "wpa2":
		intf = &OpenSecurity{}
	}

	if intf != nil {
		intf.Init()
	}

	return
}
