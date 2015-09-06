package security

type Security interface {
	Init()
	Type() string
	Params() []*SecurityParam
	Set(string, string) error
	Properties() map[string]string
}
