package security

type Security interface {
	Init()
	Type() string
	Params() []*SecurityParam
	Set(string, string) error
	Import(map[string]string)
	Export() map[string]string
}

func GetSecurity(typ string) (intf Security) {
	switch typ {
	case "open":
		intf = &OpenSecurity{}
	case "wpa", "wpa2":
		intf = &WpaSecurity{}
	}

	if intf != nil {
		intf.Init()
	}

	return
}
