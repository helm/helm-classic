package action

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/deis/helm/log"
	"github.com/deis/helm/release"
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

func TestCheckLatest(t *testing.T) {
	var oldRepo = release.Project
	var oldOwner = release.Owner
	var b bytes.Buffer
	defer func() {
		release.Project = oldRepo
		release.Owner = oldOwner
		log.Stdout = os.Stdout
	}()

	log.IsDebugging = true
	log.Stdout = &b

	// Once there is a release greater than 0.0.1, we can remove this.
	release.Project = "glide"
	release.Owner = "Masterminds"

	CheckLatest("0.0.1")

	if !strings.Contains(b.String(), "A new version of Helm") {
		t.Error("Expected notification of a new release.")
	}
}
