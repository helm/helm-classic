package release

import (
	"github.com/google/go-github/github"
)

// Owner is default Helm repository owner or organization.
var Owner = "helm"

// Project is the default Helm repository name.
var Project = "helm-classic"

// RepoService is a GitHub client instance.
var RepoService GHRepoService

// GHRepoService is a restricted interface to GitHub client operations.
type GHRepoService interface {
	GetLatestRelease(string, string) (*github.RepositoryRelease, *github.Response, error)
}

// Latest returns information on the latest Helm Classic version.
func Latest() (*github.RepositoryRelease, error) {
	if RepoService == nil {
		RepoService = github.NewClient(nil).Repositories
	}
	rel, _, err := RepoService.GetLatestRelease(Owner, Project)
	return rel, err
}

// LatestVersion returns the string version for the latest release.
func LatestVersion() (string, error) {
	rel, err := Latest()
	if err != nil {
		return "", err
	}

	if rel.TagName == nil {
		return "", nil
	}

	return *rel.TagName, nil
}

// LatestDownloadURL returns the URL from which to download a release.
func LatestDownloadURL() (string, error) {
	src, err := Latest()
	if err != nil || src.HTMLURL == nil {
		return "", err
	}
	return *src.HTMLURL, err
}
