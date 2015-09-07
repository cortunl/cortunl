package security

import (
	"github.com/dropbox/godropbox/errors"
)

type OpenSecurity struct {
	props map[string]string
}

func (s *OpenSecurity) Init() {
	s.props = map[string]string{}
	s.props["enctype"] = "open" // TODO
}

func (s *OpenSecurity) Type() string {
	return "open"
}

func (s *OpenSecurity) Params() (params []*SecurityParam) {
	return
}

func (s *OpenSecurity) Set(key, val string) (err error) {
	err = &UnknownError{
		errors.New("security.wpa: Unknown security parameter"),
	}
	return
}

func (s *OpenSecurity) Properties() map[string]string {
	return s.props
}
