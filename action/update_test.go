package action

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/helm/helm/log"
	"github.com/helm/helm/release"
)

func TestCheckLatest(t *testing.T) {
	setupTestCheckLatest()
	var b bytes.Buffer
	defer func() {
		log.Stdout = os.Stdout
		log.Stderr = os.Stderr
		release.RepoService = nil
	}()

	log.IsDebugging = true
	log.Stdout = &b
	log.Stderr = &b

	CheckLatest("0.0.1")

	if !strings.Contains(b.String(), "A new version of Helm") {
		t.Error("Expected notification of a new release.")
	}
}

type MockGHRepoService struct {
	Release *github.RepositoryRelease
}

func setupTestCheckLatest() {
	v := "9.8.7"
	u := "http://example.com/latest/release"
	i := 987
	r := &github.RepositoryRelease{
		TagName: &v,
		HTMLURL: &u,
		ID:      &i,
	}
	release.RepoService = &MockGHRepoService{Release: r}
}

func (m *MockGHRepoService) GetLatestRelease(o, p string) (*github.RepositoryRelease, *github.Response, error) {
	return m.Release, nil, nil
}
