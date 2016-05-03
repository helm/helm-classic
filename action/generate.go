package action

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/helm/helm-classic/generator"
	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/util"
)

// Generate runs generators on the entire chart.
//
// By design, this only operates on workspaces, as it should never be run
// on the cache.
func Generate(chart, homedir string, exclude []string, force bool) {
	if abs, err := filepath.Abs(homedir); err == nil {
		homedir = abs
	}
	chartPath := util.WorkspaceChartDirectory(homedir, chart)

	// Although helmc itself may use the new HELMC_HOME environment variable to optionally define its
	// home directory, to maintain compatibility with charts created for the ORIGINAL helm, we
	// continue to support expansion of these "legacy" environment variables, including HELM_HOME.
	os.Setenv("HELM_HOME", homedir)
	os.Setenv("HELM_DEFAULT_REPO", mustConfig(homedir).Repos.Default)
	os.Setenv("HELM_FORCE_FLAG", strconv.FormatBool(force))

	count, err := generator.Walk(chartPath, exclude, force)
	if err != nil {
		log.Die("Failed to complete generation: %s", err)
	}
	log.Info("Ran %d generators.", count)
}
