package release

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestLatest(t *testing.T) {
	rr, err := Latest()
	if err != nil {
		t.Errorf("Failed to get latest: %s", err)
	}

	if *rr.ID != 987 {
		t.Errorf("ID below zero.")
	}
}

func TestLatestVersion(t *testing.T) {
	v, err := LatestVersion()
	if err != nil {
		t.Error(err)
	}

	if v != "9.8.7" {
		t.Error("Expected tag, not empty string")
	}
}

func TestLatestDownloadURL(t *testing.T) {
	v, err := LatestDownloadURL()
	if err != nil {
		t.Error(err)
	}

	if v != "http://example.com/latest/release" {
		t.Error("Expected URL, not empty string")
	}

}

type MockGHRepoService struct {
	Release *github.RepositoryRelease
}

func init() {
	v := "9.8.7"
	u := "http://example.com/latest/release"
	i := 987
	r := &github.RepositoryRelease{
		TagName: &v,
		HTMLURL: &u,
		ID:      &i,
	}
	RepoService = &MockGHRepoService{Release: r}
}

func (m *MockGHRepoService) GetLatestRelease(o, p string) (*github.RepositoryRelease, *github.Response, error) {
	return m.Release, nil, nil
}
