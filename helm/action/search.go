package action

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/deis/helm/helm/log"
	"github.com/deis/helm/helm/model"
)

// Search looks for packages with 'term' in their name.
func Search(term, homedir string) {

	dirs, err := search(term, homedir)
	if err != nil {
		log.Die(err.Error())
	}

	log.Info("\n=================")
	log.Info("Available Charts")
	log.Info("=================\n")

	for _, d := range dirs {
		y, err := model.LoadChartfile(filepath.Join(d, "Chart.yaml"))
		if err != nil {
			log.Info("\t%s - UNKNOWN", filepath.Base(d))
			continue
		}
		log.Info("\t%s (%s %s) - %s", filepath.Base(d), y.Name, y.Version, y.Description)
	}

	log.Info("")
}

func search(term, homedir string) ([]string, error) {
	term = sanitizeTerm(term)
	sp := filepath.Join(homedir, CacheChartPath, term)
	dirs, err := filepath.Glob(sp)
	if err != nil {
		return dirs, fmt.Errorf("No results found. %s", err)
	} else if len(dirs) == 0 {
		return dirs, errors.New("No results found.")
	}
	return dirs, nil
}

func sanitizeTerm(term string) string {
	if term == "" {
		term = "*"
	}

	if term != "*" {
		term = "*" + term + "*"
	}

	return term
}
