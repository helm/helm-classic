// Package dependency provides dependency resolution.
package dependency

import (
	"os"
	"path/filepath"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/config"
	"github.com/helm/helm/log"
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
func Resolve(cf *chart.Chartfile, installdir string) ([]*chart.Dependency, error) {
	if len(cf.Dependencies) == 0 {
		log.Debug("No dependencies to check. :achievement-unlocked:")
		return []*chart.Dependency{}, nil
	}

	cache, err := dependencyCache(installdir)
	if err != nil {
		log.Debug("Failed to build dependency cache: %s", err)
		return []*chart.Dependency{}, err
	}

	res := []*chart.Dependency{}

	// TODO: This could be made more efficient.
	for _, check := range cf.Dependencies {
		resolved := false
		for n, chart := range cache {
			log.Debug("Checking if %s (%s) %s meets %s %s", chart.Name, n, chart.Version, check.Name, check.Version)
			if chart.From != nil {
				log.Debug("✔︎")
				if satisfies(chart.From, check) {
					resolved = true
					break
				}
			} else {
				log.Info("Chart %s is pre-0.2.0. Legacy mode enabled.", chart.Name)
				if chart.Name == check.Name && check.VersionOK(chart.Version) {
					log.Debug("✔︎")
					resolved = true
					break
				}
			}
		}
		if !resolved {
			log.Debug("No matches found for %s %s", check.Name, check.Version)
			res = append(res, check)
		}
	}
	return res, nil
}

// satisfies checks that this satisfies the dependency spec in that.
func satisfies(this, that *chart.Dependency) bool {
	if this.Name != that.Name {
		return false
	}
	if !optRepoMatch(this, that) {
		return false
	}
	return that.VersionOK(this.Version)
}

func optRepoMatch(from, req *chart.Dependency) bool {
	// If no repo is set, this is treated as a match.
	if req.Repo == "" {
		return true
	}
	// Some day we might want to do some git-fu to match different forms of the
	// same Git repo.
	a, err := canonicalRepo(req.Repo)
	if err != nil {
		log.Err("Could not parse %s: %s", req.Repo, err)
		return false
	}
	b, err := canonicalRepo(from.Repo)
	if err != nil {
		log.Err("Could not parse %s: %s", from.Repo, err)
		return false
	}
	return a == b
}

// reposMatch compares two repository URLs and returns true if they point to the same repo.
//
// This canonicalizes both repository URLs and then compares. If a repository URL
// fails to parse, this returns false.
func reposMatch(a, b string) bool {
	aa, err := canonicalRepo(a)
	if err != nil {
		log.Debug("Cannot format %s (a)", a)
		return false
	}

	bb, err := canonicalRepo(b)
	if err != nil {
		log.Debug("Cannot format %s (b)", b)
		return false
	}

	return aa == bb
}

// dependencyCache builds a map of chart and Chartfile.
func dependencyCache(chartdir string) (map[string]*chart.Chartfile, error) {
	cache := map[string]*chart.Chartfile{}
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
		cf, err := chart.LoadChartfile(filepath.Join(chartdir, fi.Name(), "Chart.yaml"))
		if err != nil {
			// If the chartfile does not load, we ignore it.
			continue
		}

		cache[fi.Name()] = cf
	}
	return cache, nil
}

// canonicalRepo is deprecated. Use dependency.CanonicalRepo instead.
func canonicalRepo(name string) (string, error) {
	return config.CanonicalRepo(name)
}
