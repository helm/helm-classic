package action_test

import (
	"path/filepath"
	"testing"

	"github.com/deis/helm/helm/action"
)

var HOME = ""

func init() {
	HOME, _ = filepath.Abs("../testdata/helm_home")
}

func TestInfo(t *testing.T) {
	action.Info("alpine", HOME)
	//TODO: assert results
}
