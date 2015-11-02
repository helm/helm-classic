package release

import (
	"github.com/google/go-github/github"
)

const owner = "deis"
const project = "helm"

// GHClient is a GitHub client.
var GHClient = github.NewClient(nil)

// Latest returns information on the latest Helm version.
func Latest() (*github.RepositoryRelease, error) {
	rel, _, err := GHClient.Repositories.GetLatestRelease(owner, project)
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

// LatestDowloadURL returns the URL to download a release.
func LatestDownloadURL() (string, error) {
	src, err := Latest()
	if err != nil || src.HTMLURL == nil {
		return "", err
	}
	return *src.HTMLURL, err
}
