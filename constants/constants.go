package constants

import (
	"path/filepath"
)

const (
	Name         = "cortunl"
	NameFormated = "Cortunl"
	Version      = "0.1.0"
)

var (
	TempDir = filepath.Join(string(filepath.Separator), "tmp", "cortunl")
)
