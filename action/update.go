package action

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Masterminds/vcs"

	"github.com/deis/helm/log"
	"github.com/deis/helm/repo"
)

// Update fetches the remote repo into the home directory.
func Update(home string) {
	home, err := filepath.Abs(home)
	if err != nil {
		log.Die("Could not generate absolute path for %q: %s", home, err)
	}

	// Basically, install if this is the first run.
	ensurePrereqs()
	ensureHome(home)

	rpath := filepath.Join(home, CachePath, Repofile)
	rc, err := repo.LoadRepofile(rpath)
	if err != nil {
		log.Die(err.Error())
	}
	if err := rc.UpdateAll(); err != nil {
		log.Die("Not all repos could be updated: %s", err)
	}
	log.Info("Done")
}

// gitUpdate updates a Git repo.
func gitUpdate(git *vcs.GitRepo) error {
	if err := git.Update(); err != nil {
		return err
	}

	log.Debug("Updated %s from %s", git.LocalPath(), git.Remote())
	return nil
}

// ensurePrereqs verifies that Git and Kubectl are both available.
func ensurePrereqs() {
	if _, err := exec.LookPath("git"); err != nil {
		log.Die("Could not find 'git' on $PATH: %s", err)
	}
	if _, err := exec.LookPath("kubectl"); err != nil {
		log.Die("Could not find 'kubectl' on $PATH: %s", err)
	}
}

// ensureRepo ensures that the repo exists and is checked out.
// DEPRECATED: You should use the functions in package `repo` instead.
func ensureRepo(repo, home string) *vcs.GitRepo {
	if err := os.Chdir(home); err != nil {
		log.Die("Could not change to directory %q: %s", home, err)
	}
	git, err := vcs.NewGitRepo(repo, home)
	if err != nil {
		log.Die("Could not get repository %q: %s", repo, err)
	}

	git.Logger = log.New()

	if !git.CheckLocal() {
		log.Debug("Cloning repo into %q. Please wait.", home)
		if err := git.Get(); err != nil {
			log.Die("Could not create repository in %q: %s", home, err)
		}
	}

	return git
}

// ensureHome ensures that a HELM_HOME exists.
func ensureHome(home string) {

	must := []string{home, filepath.Join(home, CachePath), filepath.Join(home, WorkspacePath), filepath.Join(home, CacheChartPath)}

	for _, p := range must {
		if fi, err := os.Stat(p); err != nil {
			log.Debug("Creating %s", p)
			if err := os.MkdirAll(p, 0755); err != nil {
				log.Die("Could not create %q: %s", p, err)
			}
		} else if !fi.IsDir() {
			log.Die("%s must be a directory.", home)
		}
	}

	refi := filepath.Join(home, CachePath, Repofile)
	if _, err := os.Stat(refi); err != nil {
		log.Info("Creating %s", refi)
		// Attempt to create a Repos.yaml
		if err := ioutil.WriteFile(refi, []byte(repo.DefaultRepofile), 0755); err != nil {
			log.Die("Could not create %s: %s", refi, err)
		}
	}

	if err := os.Chdir(home); err != nil {
		log.Die("Could not change to directory %q: %s", home, err)
	}
}
