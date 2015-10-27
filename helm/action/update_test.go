package action

import (
	"os"
	"path/filepath"
	"testing"
)

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
	repo := "https://github.com/deis/helm"
	ensureRepo(repo, filepath.Join(HOME, "cache"))
}
