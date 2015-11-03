package action

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/deis/helm/log"
	"github.com/deis/helm/model"
)

// Search looks for packages with 'term' in their name.
func Search(term, homedir string) {
	charts, err := search(term, homedir)
	if err != nil {
		log.Die(err.Error())
	}

	log.Info("\n=================")
	log.Info("Available Charts")
	log.Info("=================\n")

	log.Info("")

	for dir, chart := range charts {
		log.Info("\t%s (%s %s) - %s", filepath.Base(dir), chart.Name, chart.Version, chart.Description)
	}
}

func search(term, homedir string) (map[string]*model.Chartfile, error) {
	dirs, err := filepath.Glob(filepath.Join(homedir, CacheChartPath, "*"))

	if err != nil {
		return nil, fmt.Errorf("No results found. %s", err)
	} else if len(dirs) == 0 {
		return nil, errors.New("No results found.")
	}

	charts := make(map[string]*model.Chartfile)

	r, _ := regexp.Compile(term)

	for _, dir := range dirs {
		chart, err := model.LoadChartfile(filepath.Join(dir, "Chart.yaml"))

		if err != nil {
			log.Info("\t%s - UNKNOWN", filepath.Base(dir))
			continue
		} else if r.MatchString(chart.Name) || r.MatchString(chart.Description) {
			charts[dir] = chart
		}
	}

	return charts, nil
}
