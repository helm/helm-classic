package action

import (
	"os"
	"path"
)

import (
	"github.com/deis/helm/log"
)

// Publish a chart from the workspace to the cache directory
//
// - chartName being published
// - homeDir is the helm home directory for the user
// - force publishing even if the chart directory already exists
func Publish(chartName, homeDir string, force bool) {

	src := path.Join(homeDir, WorkspaceChartPath, chartName)
	dst := path.Join(homeDir, CacheChartPath, chartName)

	if _, err := os.Stat(dst); err == nil {
		if force != true {
			log.Info("chart already exists, use -f to force")
			return
		}
	}

	if err := copyDir(src, dst); err != nil {
		log.Die("failed to publish directory: %v", err)
	}
}
