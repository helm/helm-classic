package action

import (
	"path/filepath"

	"github.com/Masterminds/semver"

	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/release"
)

// Update fetches the remote repo into the home directory.
func Update(home string) {
	home, err := filepath.Abs(home)
	if err != nil {
		log.Die("Could not generate absolute path for %q: %s", home, err)
	}

	CheckLocalPrereqs(home)

	rc := mustConfig(home).Repos
	if err := rc.UpdateAll(); err != nil {
		log.Die("Not all repos could be updated: %s", err)
	}
	log.Info("Done")
}

// CheckLatest checks whether this version of Helm Classic is the latest version.
//
// This does not ensure that this is the latest. If a newer version is found,
// this generates a message indicating that.
//
// The passed-in version is the base version that will be checked against the
// remote release list.
func CheckLatest(version string) {
	ver, err := release.LatestVersion()
	if err != nil {
		log.Warn("Skipped Helm Classic version check: %s", err)
		return
	}

	current, err := semver.NewVersion(version)
	if err != nil {
		log.Warn("Local version %s is not well-formed", version)
		return
	}
	remote, err := semver.NewVersion(ver)
	if err != nil {
		log.Warn("Remote version %s is not well-formed", ver)
		return
	}

	if remote.GreaterThan(current) {
		log.Warn("A new version of Helm Classic is available. You have %s. The latest is %v", version, ver)
		log.Info("Download version %s by running: %s", ver, "curl -s https://get.helm.sh | bash")
	}

}
