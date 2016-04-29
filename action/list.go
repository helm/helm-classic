package action

import (
	"path/filepath"

	"github.com/helm/helm-classic/chart"
	"github.com/helm/helm-classic/log"
	helm "github.com/helm/helm-classic/util"
)

// List lists all of the local charts.
func List(homedir string) {
	md := helm.WorkspaceChartDirectory(homedir, "*")
	charts, err := filepath.Glob(md)
	if err != nil {
		log.Warn("Could not find any charts in %q: %s", md, err)
	}
	for _, c := range charts {
		cname := filepath.Base(c)
		if ch, err := chart.LoadChartfile(filepath.Join(c, Chartfile)); err == nil {
			log.Info("\t%s (%s %s) - %s", cname, ch.Name, ch.Version, ch.Description)
			continue
		}
		log.Info("\t%s (unknown)", cname)
	}
}
