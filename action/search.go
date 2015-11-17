package action

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/log"
)

// Search looks for packages with 'term' in their name.
func Search(term, homedir string) {
	charts, err := searchAll(term, homedir)
	if err != nil {
		log.Die(err.Error())
	}

	if len(charts) > 0 {
		for _, name := range sortedIndex(charts) {
			chart := charts[name]
			log.Msg("\t%s (%s %s) - %s", name, chart.Name, chart.Version, chart.Description)
		}
	} else {
		log.Msg("No chart found for \"%s\".", term)
	}

}

func sortedIndex(m map[string]*chart.Chartfile) []string {
	ss := make(sort.StringSlice, len(m))

	i := 0
	for k := range m {
		ss[i] = k
		i++
	}

	ss.Sort()
	return ss
}

func searchAll(term, homedir string) (map[string]*chart.Chartfile, error) {
	r := mustConfig(homedir).Repos
	results := map[string]*chart.Chartfile{}
	for _, table := range r.Tables {
		tablename := table.Name
		if table.Name == r.Default {
			tablename = ""
		}
		base := filepath.Join(homedir, CachePath, table.Name, "*")
		if err := search(term, base, tablename, results); err != nil {
			log.Warn("Search error: %s", err)
		}
	}
	return results, nil
}

func search(term, base, table string, charts map[string]*chart.Chartfile) error {
	dirs, err := filepath.Glob(base)
	if err != nil {
		return fmt.Errorf("No results found. %s", err)
	} else if len(dirs) == 0 {
		return errors.New("No results found.")
	}

	r, err := regexp.Compile(term)
	if err != nil {
		log.Die("Invalid expression %q: %s", term, err)
	}

	for _, dir := range dirs {
		cname := filepath.Join(table, filepath.Base(dir))
		chrt, err := chart.LoadChartfile(filepath.Join(dir, "Chart.yaml"))

		if err != nil {
			// This dir is not a chart. Skip it.
			continue
		} else if r.MatchString(chrt.Name) || r.MatchString(chrt.Description) {
			charts[cname] = chrt
		}
	}

	return nil
}
