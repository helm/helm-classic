package action

import (
	"os"
	"path/filepath"
	"testing"
)

var HOME = ""

func init() {
	HOME = filepath.Join(helmRoot, "helm/testdata/helm_home")
}

func TestEnsurePrereqs(t *testing.T) {
	pp := os.Getenv("PATH")
	defer os.Setenv("PATH", pp)
	here, _ := os.Getwd()
	os.Setenv("PATH", filepath.Join(here, "../testdata")+":"+pp)
	ensurePrereqs()
}

func TestEnsureHome(t *testing.T) {
	ensureHome(HOME)
}

func TestEnsureRepo(t *testing.T) {
	repo := "https://github.com/deis/charts"
	ensureRepo(repo, filepath.Join(HOME, "cache", "charts"))
}
