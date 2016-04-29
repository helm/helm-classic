package action

import "os"

import (
	"github.com/helm/helm-classic/log"
	helm "github.com/helm/helm-classic/util"
)

// Publish a chart from the workspace to the cache directory
//
// - chartName being published
// - homeDir is the helm home directory for the user
// - force publishing even if the chart directory already exists
func Publish(chartName, homeDir, repo string, force bool) {
	if repo == "" {
		repo = "charts"
	}

	if !mustConfig(homeDir).Repos.Exists(repo) {
		log.Err("Repo %s does not exist", repo)
		log.Info("Available repositories")
		ListRepos(homeDir)
		return
	}

	src := helm.WorkspaceChartDirectory(homeDir, chartName)
	dst := helm.CacheDirectory(homeDir, repo, chartName)

	if _, err := os.Stat(dst); err == nil {
		if force != true {
			log.Info("chart already exists, use -f to force")
			return
		}
	}

	if err := helm.CopyDir(src, dst); err != nil {
		log.Die("failed to publish directory: %v", err)
	}
}
