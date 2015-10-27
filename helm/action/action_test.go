package action

import (
	"path/filepath"
)

var HOME = ""

func init() {
	HOME, _ = filepath.Abs("../testdata/helm_home")
}
