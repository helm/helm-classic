package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/helm/helm-classic/log"
)

// Plugin attempts to execute a plugin.
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
// 	- $HELM_COMMAND: the name of the command (as seen by Helm) that resulted in this program being executed.
func Plugin(homedir, cmd string, args []string) {
	if abs, err := filepath.Abs(homedir); err == nil {
		homedir = abs
	}

	// Although helmc itself may use the new HELMC_HOME environment variable to optionally define its
	// home directory, to maintain compatibility with plugins created for the ORIGINAL helm, we
	// continue to support expansion of these "legacy" environment variables, including HELM_HOME.
	os.Setenv("HELM_HOME", homedir)
	os.Setenv("HELM_COMMAND", args[0])
	os.Setenv("HELM_DEFAULT_REPO", mustConfig(homedir).Repos.Default)

	cmd = PluginName(cmd)
	execPlugin(cmd, args[1:])
}

// HasPlugin returns true if the named plugin exists.
func HasPlugin(name string) bool {
	name = PluginName(name)
	_, err := exec.LookPath(name)
	return err == nil
}

// PluginName returns the full plugin name.
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
