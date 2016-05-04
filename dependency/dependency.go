// Package dependency provides dependency resolution.
package dependency

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/helm/helm-classic/chart"
	"github.com/helm/helm-classic/log"
)

// Resolve takes a chart and a location and checks whether the chart's dependencies are satisfied.
//
// The `installdir` is the location where installed charts are located. Typically
// this is in $HELMC_HOME/workspace/charts.
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

// canonicalRepo returns a canonical repo name of the form `host/path.git`.
//
// There are several accepted Git protocol representations:
//
//	- /PATH.git (local)   -> localhost/PATH.git
//	- file:///PATH.git (local)   -> localhost/PATH.git
//	- https://HOST/PATH.git    -> HOST/PATH.git
//	- http://HOST/PATH.git    -> HOST/PATH.git
//	- ssh://user@HOST/PATH.git -> HOST/PATH.git
//	- user@HOST:PATH.git  -> HOST/PATH.git
//
// In the case where no suitable normalization can be found, this will return
// the original string, assuming that there is some additional Git representation
// that we don't know about.
func canonicalRepo(name string) (string, error) {

	if strings.Index(name, "://") > 0 {
		// URL parseable
		u, err := url.Parse(name)
		if err != nil {
			return name, err
		}

		if u.Scheme == "file" && u.Host == "" {
			u.Host = "localhost"
		}

		return filepath.Join(u.Host, u.Path), nil
	} else if i := strings.Index(name, "@"); i > 0 && i < strings.Index(name, ":") {

		a := strings.SplitN(name, "@", 2)
		if len(a) != 2 {
			return name, fmt.Errorf("Could not parse SCP name %s: '@' split failed", name)
		}

		a = strings.SplitN(a[1], ":", 2)
		if len(a) != 2 {
			return name, fmt.Errorf("Could not parse SCP name %s: ':' split failed", name)
		}

		return filepath.Join(a[0], a[1]), nil
	}
	// Is a filepath
	return filepath.Join("localhost", name), nil
}
