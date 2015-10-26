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
	// Skip right now. This is covered in issue #58, and fixed in the associated
	// PR.
	t.Skip()
	action.Info("alpine", HOME)
	//TODO: assert results
}
