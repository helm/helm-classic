package action

import (
	"os/exec"

	"github.com/helm/helm-classic/log"
	helm "github.com/helm/helm-classic/util"
)

// Doctor helps you see what's wrong with your Helm Classic setup
func Doctor(home string) {
	log.Info("Checking things locally...")
	CheckLocalPrereqs(home)
	CheckKubePrereqs()

	log.Info("Everything looks good! Happy helming!")
}

// CheckAllPrereqs makes sure we have all the tools we need for overall
// Helm Classic success
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
	helm.EnsureHome(home)
	ensureCommand("git")
}

func ensureCommand(command string) {
	if _, err := exec.LookPath(command); err != nil {
		log.Die("Could not find '%s' on $PATH: %s", command, err)
	}
}
