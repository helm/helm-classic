package action

import (
	"path/filepath"

	"github.com/helm/helm-classic/config"
	"github.com/helm/helm-classic/log"
	helm "github.com/helm/helm-classic/util"
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
		log.Warn("Oops! Looks like we had some issues running your command! Running `helmc doctor` to ensure we have all the necessary prerequisites in place...")
		Doctor(homedir)
		cfg, err = config.Load(rpath)
		if err != nil {
			log.Die("Oops! Could not load %s. Error: %s", rpath, err)
		}
		log.Info("Continuing onwards and upwards!")
	}
	return cfg
}
