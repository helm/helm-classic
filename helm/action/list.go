package action

import (
	"path/filepath"

	"github.com/deis/helm/helm/model"
)

// List lists all of the local charts.
func List(homedir string) {
	md := filepath.Join(homedir, WorkdirChartPath, "*")
	charts, err := filepath.Glob(md)
	if err != nil {
		Warn("Could not find any charts in %q: %s", md, err)
	}
	for _, c := range charts {
		cname := filepath.Base(c)
		if ch, err := model.Load(filepath.Join(c, "Chart.yaml")); err == nil {
			Info("\t%s (%s %s) - %s", cname, ch.Name, ch.Version, ch.Description)
			continue
		}
		Info("\t%s (unknown)", cname)
	}
}
