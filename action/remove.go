package action

import (
	"os"
	"path/filepath"

	"github.com/helm/helm/log"
)

// Remove removes a chart from the workdir.
//
// - chart is the source
// - homedir is the home directory for the user
func Remove(chart string, homedir string) {
	chartPath := filepath.Join(homedir, WorkspaceChartPath, chart)
	if _, err := os.Stat(chartPath); err != nil {
		log.Die("Chart not found. %s", err)
	}

	if err := os.RemoveAll(chartPath); err != nil {
		log.Die("%s", err)
	}

	log.Info("All clear! You have successfully removed %s from your workspace.", chart)
}
