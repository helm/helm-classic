package action

import (
	"os"
	"os/exec"

	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/util"
)

// Edit charts using the shell-defined $EDITOR
//
// - chartName being edited
// - homeDir is the Helm Classic home directory for the user
func Edit(chartName, homeDir string) {

	chartDir := util.WorkspaceChartDirectory(homeDir, chartName)

	if _, err := os.Stat(chartDir); os.IsNotExist(err) {
		log.Die("Could not find chart: %s", chartName)
	}

	openEditor(chartDir)
}

// openEditor opens the given filename in an interactive editor
func openEditor(path string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		log.Die("must set shell $EDITOR")
	}

	editorPath, err := exec.LookPath(editor)
	if err != nil {
		log.Die("Could not find %s in PATH", editor)
	}

	cmd := exec.Command(editorPath, path)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		log.Die("Could not open $EDITOR: %s", err)
	}
}
