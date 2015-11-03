package action

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/deis/helm/chart"
	"github.com/deis/helm/log"
)

// Search looks for packages with 'term' in their name.
func Search(term, homedir string) {
	charts, err := search(term, homedir)
	if err != nil {
		log.Die(err.Error())
	}

	log.Info("=================")
	log.Info("Available Charts")
	log.Info("=================")

	log.Info("")

	for dir, chart := range charts {
		log.Info("\t%s (%s %s) - %s", filepath.Base(dir), chart.Name, chart.Version, chart.Description)
	}
}

func search(term, homedir string) (map[string]*chart.Chartfile, error) {
	files, err := filepath.Glob(filepath.Join(homedir, CacheChartPath, "*"))

	// only return chart directories
	var dirs []string
	for _, f := range files {
		if filepath.Base(f) == ".git" {
			continue
		}
		fm, _ := os.Stat(f)
		if !fm.IsDir() {
			continue
		}
		dirs = append(dirs, f)
	}

	if err != nil {
		return nil, fmt.Errorf("No results found. %s", err)
	} else if len(dirs) == 0 {
		return nil, errors.New("No results found.")
	}

	charts := make(map[string]*chart.Chartfile)

	r, _ := regexp.Compile(term)

	for _, dir := range dirs {
		chart, err := chart.LoadChartfile(filepath.Join(dir, "Chart.yaml"))

		if err != nil {
			log.Warn("failed to load Chart.yaml: %v", err)
			continue
		} else if r.MatchString(chart.Name) || r.MatchString(chart.Description) {
			charts[dir] = chart
		}
	}

	return charts, nil
}
