package security

import (
	"github.com/dropbox/godropbox/errors"
)

type OpenSecurity struct {
	props map[string]string
}

func (s *OpenSecurity) Init() {
	s.props = map[string]string{}
}

func (s *OpenSecurity) Type() string {
	return "open"
}

func (s *OpenSecurity) Params() (params []*SecurityParam) {
	return
}

func (s *OpenSecurity) Set(key, val string) (err error) {
	err = &UnknownError{
		errors.New("security.open: Unknown security parameter"),
	}
	return
}

func (s *OpenSecurity) Import(props map[string]string) {
	s.props = props
}

func (s *OpenSecurity) Export() map[string]string {
	return s.props
}
