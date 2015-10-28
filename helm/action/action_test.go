package action

import (
	"path/filepath"
)

var TestHome string

func init() {
	TestHome, _ = filepath.Abs("../testdata/helm_home")
}
