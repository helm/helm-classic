package action

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helm/helm-classic/test"
	"github.com/helm/helm-classic/util"
)

func TestEnsurePrereqs(t *testing.T) {
	pp := os.Getenv("PATH")
	defer os.Setenv("PATH", pp)

	os.Setenv("PATH", filepath.Join(test.HelmRoot, "testdata")+":"+pp)

	homedir := test.CreateTmpHome()
	CheckAllPrereqs(homedir)
}

func TestEnsureHome(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	util.EnsureHome(tmpHome)
}
