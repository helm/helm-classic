package action

import (
	"os"
	"path"
)

import (
	"github.com/helm/helm/log"
	helm "github.com/helm/helm/util"
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

	src := path.Join(homeDir, helm.WorkspaceChartPath, chartName)
	dst := path.Join(homeDir, helm.CachePath, repo, chartName)

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
