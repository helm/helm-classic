/* Package dependency provides dependency resolution.*/
package dependency

import (
	"os"
	"path/filepath"

	"github.com/deis/helm/log"
	"github.com/deis/helm/model"
)

// Resolve takes a chart and a location and checks whether the chart's dependencies are satisfied.
//
// The `installdir` is the location where installed charts are located. Typically
// this is in $HELM_HOME/workspace/charts.
//
// This returns a list of unsatisfied dependencies (NOT an error condition).
//
// It returns an error only if it cannot perform the task of resolving dependencies.
// Failed dependencies to not constitute an error.
func Resolve(cf *model.Chartfile, installdir string) ([]*model.Dependency, error) {
	if len(cf.Dependencies) == 0 {
		log.Debug("No dependencies to check. :achievement-unlocked:")
		return []*model.Dependency{}, nil
	}

	cache, err := dependencyCache(installdir)
	if err != nil {
		log.Debug("Failed to build dependency cache: %s", err)
		return []*model.Dependency{}, err
	}

	res := []*model.Dependency{}

	// TODO: This could be made more efficient.
	for _, check := range cf.Dependencies {
		resolved := false
		for n, chart := range cache {
			log.Debug("Checking if %s (%s) %s meets %s %s", chart.Name, n, chart.Version, check.Name, check.Version)
			if chart.Name == check.Name && check.VersionOK(chart.Version) {
				log.Debug("✔︎")
				resolved = true
				break
			}
		}
		if !resolved {
			log.Debug("No matches found for %s %s", check.Name, check.Version)
			res = append(res, check)
		}
	}
	return res, nil
}

// dependencyCache builds a map of chart and Chartfile.
func dependencyCache(chartdir string) (map[string]*model.Chartfile, error) {
	cache := map[string]*model.Chartfile{}
	dir, err := os.Open(chartdir)
	if err != nil {
		return cache, err
	}
	defer dir.Close()

	fis, err := dir.Readdir(0)
	if err != nil {
		return cache, err
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		cf, err := model.LoadChartfile(filepath.Join(chartdir, fi.Name(), "Chart.yaml"))
		if err != nil {
			// If the chartfile does not load, we ignore it.
			continue
		}

		cache[fi.Name()] = cf
	}
	return cache, nil
}
