package action

import (
	"os"
	"path/filepath"

	"github.com/helm/helm/generator"
	"github.com/helm/helm/log"
	"github.com/helm/helm/util"
)

// Generate runs generators on the entire chart.
//
// By design, this only operates on workspaces, as it should never be run
// on the cache.
func Generate(chart, homedir string, exclude []string) {
	if abs, err := filepath.Abs(homedir); err == nil {
		homedir = abs
	}
	chartPath := util.WorkspaceChartDirectory(homedir, chart)

	os.Setenv("HELM_HOME", homedir)
	os.Setenv("HELM_DEFAULT_REPO", mustConfig(homedir).Repos.Default)

	count, err := generator.Walk(chartPath, exclude)
	if err != nil {
		log.Die("Failed to complete generation: %s", err)
	}
	log.Info("Ran %d generators.", count)
}
