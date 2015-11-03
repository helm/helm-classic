package action

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsurePrereqs(t *testing.T) {
	pp := os.Getenv("PATH")
	defer os.Setenv("PATH", pp)

	os.Setenv("PATH", filepath.Join(helmRoot, "testdata")+":"+pp)
	ensurePrereqs()
}

func TestEnsureHome(t *testing.T) {
	tmpHome := createTmpHome()
	ensureHome(tmpHome)
}

func TestEnsureRepo(t *testing.T) {
	tmpHome := createTmpHome()
	ensureHome(tmpHome)

	repo := "https://github.com/deis/charts"
	ensureRepo(repo, filepath.Join(tmpHome, "cache", "charts"))
}
