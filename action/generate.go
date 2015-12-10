package action

import (
	"os"
	"path/filepath"

	"github.com/helm/helm/generator"
	"github.com/helm/helm/log"
	"github.com/helm/helm/util"
)

func Generate(chart, homedir string) {
	if abs, err := filepath.Abs(homedir); err == nil {
		homedir = abs
	}
	chartPath := util.WorkspaceChartDirectory(homedir, chart)

	os.Setenv("HELM_HOME", homedir)
	os.Setenv("HELM_DEFAULT_REPO", mustConfig(homedir).Repos.Default)

	count, err := generator.Walk(chartPath)
	if err != nil {
		log.Die("Failed to complete generation: %s", err)
	}
	log.Info("Ran %d generators.", count)
}
