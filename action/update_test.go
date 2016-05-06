package action

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/release"
	"github.com/helm/helm-classic/test"
)

func TestCheckLatest(t *testing.T) {
	setupTestCheckLatest()
	defer func() {
		release.RepoService = nil
	}()

	log.IsDebugging = true

	actual := test.CaptureOutput(func() {
		CheckLatest("0.0.1")
	})

	test.ExpectContains(t, actual, "A new version of Helm Classic")
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
