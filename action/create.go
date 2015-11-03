package action

import (
	"os"
	"path"
)

import (
	"github.com/deis/helm/log"
)

// Create a chart
//
// - chartName being created
// - homeDir is the helm home directory for the user
func Create(chartName, homeDir string) {

	skeletonDir := path.Join(homeDir, CachePath, "skel")

	if fi, err := os.Stat(skeletonDir); err != nil {
		log.Die("Could not find %s: %s", skeletonDir, err)
	} else if !fi.IsDir() {
		log.Die("Malformed skeleton: %s: Must be a directory.", skeletonDir)
	}

	chartDir := path.Join(homeDir, WorkspaceChartPath, chartName)

	// copy skeleton to chart directory
	if err := copyDir(skeletonDir, chartDir); err != nil {
		log.Die("failed to copy skeleton directory: %v", err)
	}

}
