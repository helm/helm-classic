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

	homedir := createTmpHome()
	CheckAllPrereqs(homedir)
}

func TestEnsureHome(t *testing.T) {
	tmpHome := createTmpHome()
	ensureHome(tmpHome)
}
