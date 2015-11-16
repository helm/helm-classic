package action

import (
	"path/filepath"

	"github.com/helm/helm/config"
	"github.com/helm/helm/log"
)

// CachePath is the suffix for the cache.
const CachePath = "cache"

// CacheChartPath is the directory that contains a user's cached charts.
const CacheChartPath = "cache/charts"

// WorkspacePath is the user's workspace directory.
const WorkspacePath = "workspace"

// WorkspaceChartPath is the directory that contains a user's workspace charts.
const WorkspaceChartPath = "workspace/charts"

// Configfile is the file containing helm's YAML configuration data.
const Configfile = "config.yaml"

var helmpaths = []string{CachePath, WorkspacePath}

// mustConfig parses a config file or dies trying.
func mustConfig(homedir string) *config.Configfile {
	rpath := filepath.Join(homedir, Configfile)
	cfg, err := config.Load(rpath)
	if err != nil {
		log.Die("Could not load %s: %s", rpath, err)
	}
	return cfg
}
