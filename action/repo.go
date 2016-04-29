package action

import (
	"github.com/helm/helm-classic/log"
)

// ListRepos lists the repositories.
func ListRepos(homedir string) {
	rf := mustConfig(homedir).Repos

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
	cfg := mustConfig(homedir)

	if err := cfg.Repos.Add(name, repository); err != nil {
		log.Die(err.Error())
	}
	if err := cfg.Save(""); err != nil {
		log.Die("Could not save configuration: %s", err)
	}

	log.Info("Hooray! Successfully added the repo.")
}

// DeleteRepo deletes a repository.
func DeleteRepo(homedir, name string) {
	cfg := mustConfig(homedir)

	if err := cfg.Repos.Delete(name); err != nil {
		log.Die("Failed to delete repository: %s", err)
	}
	if err := cfg.Save(""); err != nil {
		log.Die("Deleted repo, but could not save settings: %s", err)
	}
}
