package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/deis/helm/log"
)

// Plugin attepmts to execute a plugin.
//
// It looks on the path for an executable named `helm-COMMAND`, and executes
// that, passing it all of the arguments received after the subcommand.
//
// Output is passed directly back to the user.
//
// This ensures that the following environment variables are set:
//
//	- $HELM_HOME: points to the user's Helm home directory.
// 	- $HELM_DEFAULT_REPO: the local name of the default repository.
func Plugin(homedir, cmd string, args []string) {
	if abs, err := filepath.Abs(homedir); err == nil {
		homedir = abs
	}
	os.Setenv("HELM_HOME", homedir)
	os.Setenv("HELM_DEFAULT_REPO", mustConfig(homedir).Repos.Default)

	cmd = PluginName(cmd)
	execPlugin(cmd, args)
}

func HasPlugin(name string) bool {
	name = PluginName(name)
	_, err := exec.LookPath(name)
	return err == nil
}

func PluginName(name string) string {
	return "helm-" + name
}

func execPlugin(name string, args []string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Die(err.Error())
	}
}
