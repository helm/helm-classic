package action_test

import (
	"path/filepath"
	"testing"

	"github.com/deis/helm/helm/action"
)

var HOME = ""

func init() {
	HOME, _ = filepath.Abs("../testdata/")
}

func TestInfo(t *testing.T) {
	action.Info("kitchensink", HOME)
	//TODO: assert results
}
