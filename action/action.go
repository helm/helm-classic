package action

import (
	"path/filepath"

	"github.com/deis/helm/config"
	"github.com/deis/helm/log"
)

// CachePath is the suffix for the cache.
const CachePath = "cache"
const CacheChartPath = "cache/charts"

const WorkspacePath = "workspace"
const WorkspaceChartPath = "workspace/charts"
const Configfile = "config.yaml"

const DefaultNS = "default"

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
