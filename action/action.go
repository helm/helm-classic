package action

import (
	"path/filepath"

	"github.com/helm/helm/config"
	"github.com/helm/helm/log"
	helm "github.com/helm/helm/util"
)

const (
	// Chartfile is the name of the YAML file that contains chart metadata.
	// One must exist inside the top level directory of every chart.
	Chartfile = "Chart.yaml"
)

// mustConfig parses a config file or dies trying.
func mustConfig(homedir string) *config.Configfile {
	rpath := filepath.Join(homedir, helm.Configfile)
	cfg, err := config.Load(rpath)
	if err != nil {
		log.Die("Could not load %s: %s", rpath, err)
	}
	return cfg
}
