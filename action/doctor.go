package action

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/helm/helm/config"
	"github.com/helm/helm/log"
)

// Doctor helps you see what's wrong with your helm setup
func Doctor(home string) {
	log.Info("Checking things locally...")
	CheckLocalPrereqs(home)
	CheckKubePrereqs()

	log.Info("Everything looks good! Happy helming!")
}

// CheckAllPrereqs makes sure we have all the tools we need for overall
// helm success
func CheckAllPrereqs(home string) {
	CheckLocalPrereqs(home)
	CheckKubePrereqs()
}

// CheckKubePrereqs makes sure we have the tools necessary to interact
// with a kubernetes cluster
func CheckKubePrereqs() {
	ensureCommand("kubectl")
}

// CheckLocalPrereqs makes sure we have all the tools we need to work with
// charts locally
func CheckLocalPrereqs(home string) {
	ensureHome(home)
	ensureCommand("git")
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

	refi := filepath.Join(home, Configfile)
	if _, err := os.Stat(refi); err != nil {
		log.Info("Creating %s", refi)
		// Attempt to create a Repos.yaml
		if err := ioutil.WriteFile(refi, []byte(config.DefaultConfigfile), 0755); err != nil {
			log.Die("Could not create %s: %s", refi, err)
		}
	}

	if err := os.Chdir(home); err != nil {
		log.Die("Could not change to directory %q: %s", home, err)
	}
}

func ensureCommand(command string) {
	if _, err := exec.LookPath(command); err != nil {
		log.Die("Could not find '%s' on $PATH: %s", command, err)
	}
}
