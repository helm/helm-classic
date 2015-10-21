package action

import (
	"path/filepath"

	"github.com/deis/helm/helm/model"
)

// Search looks for packages with 'term' in their name.
func Search(term, homedir string) {
	term = sanitizeTerm(term)
	sp := filepath.Join(homedir, CacheChartPath, "*"+term+"*")
	dirs, err := filepath.Glob(sp)
	if err != nil {
		Die("No results found. %s", err)
	}

	Info("\n=================")
	Info("Available Charts")
	Info("=================\n")

	for _, d := range dirs {
		y, err := model.Load(filepath.Join(d, "Chart.yaml"))
		if err != nil {
			Info("\t%s - UNKNOWN", filepath.Base(d))
			continue
		}
		Info("\t%s (%s %s) - %s", filepath.Base(d), y.Name, y.Version, y.Description)
	}

	Info("")
}

func sanitizeTerm(term string) string {
	return term
}
