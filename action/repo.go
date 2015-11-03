package action

import (
	"os"
	"path/filepath"

	"github.com/deis/helm/log"
	"github.com/deis/helm/repo"
)

// ListRepos lists the repositories.
func ListRepos(homedir string) {
	rp := filepath.Join(homedir, CachePath, Repofile)
	if _, err := os.Stat(rp); err != nil {
		log.Die("No YAML file found at %s", rp)
	}

	rf, err := repo.LoadRepofile(rp)
	if err != nil {
		log.Die("Could not read %s: %s", rp, err)
	}

	for _, t := range rf.Tables {
		n := t.Name
		if t.Name == rf.Default {
			n += "*"
		}
		log.Msg("\t%s\t%s", n, t.Repo)
	}
}

// AddRepo adds a repo to the list of repositories.
func AddRepo(homedir, name, repository string) {
	rpath := filepath.Join(homedir, CachePath)
	newpath := filepath.Join(rpath, name)
	rp := filepath.Join(rpath, Repofile)

	if _, err := os.Stat(newpath); err == nil {
		log.Die("A directory named %s already exists.", newpath)
	}

	rf, err := repo.LoadRepofile(rp)
	if err != nil {
		log.Die("Could not read %s: %s", rp, err)
	}

	if err := rf.Add(name, repository); err != nil {
		log.Die(err.Error())
	}

}

// DeleteRepo deletes a repository.
func DeleteRepo(homedir, name string) {
	rp := filepath.Join(homedir, CachePath, Repofile)

	rf, err := repo.LoadRepofile(rp)
	if err != nil {
		log.Die("Could not read %s: %s", rp, err)
	}

	if err := rf.Delete(name); err != nil {
		log.Die("Failed to delete repository: %s", err)
	}
}
