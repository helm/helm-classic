package action

import (
	"path/filepath"

	"github.com/deis/helm/chart"
	"github.com/deis/helm/log"
)

// List lists all of the local charts.
func List(homedir string) {
	md := filepath.Join(homedir, WorkspaceChartPath, "*")
	charts, err := filepath.Glob(md)
	if err != nil {
		log.Warn("Could not find any charts in %q: %s", md, err)
	}
	for _, c := range charts {
		cname := filepath.Base(c)
		if ch, err := chart.LoadChartfile(filepath.Join(c, "Chart.yaml")); err == nil {
			log.Info("\t%s (%s %s) - %s", cname, ch.Name, ch.Version, ch.Description)
			continue
		}
		log.Info("\t%s (unknown)", cname)
	}
}
