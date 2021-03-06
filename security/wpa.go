package security

import (
	"github.com/dropbox/godropbox/errors"
)

var wpaParams = []*SecurityParam{
	&SecurityParam{
		Label: "Password",
		Key:   "password",
	},
}

type WpaSecurity struct {
	props map[string]string
}

func (s *WpaSecurity) Init() {
	s.props = map[string]string{}
	s.props["Security"] = "wpa"
}

func (s *WpaSecurity) Type() string {
	return "wpa"
}

func (s *WpaSecurity) Params() []*SecurityParam {
	return wpaParams
}

func (s *WpaSecurity) Set(key, val string) (err error) {
	switch key {
	case "password":
		if val == "" {
			err = &InvalidError{
				errors.New("security.wpa: Password cannot be empty"),
			}
			return
		}
		s.props["Key"] = val
	default:
		err = &UnknownError{
			errors.New("security.wpa: Unknown security parameter"),
		}
		return
	}

	return
}

func (s *WpaSecurity) Import(props map[string]string) {
	s.props = props
}

func (s *WpaSecurity) Export() map[string]string {
	return s.props
}
