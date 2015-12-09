package dependency

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/config"
	"github.com/helm/helm/log"
)

// Resolver resolves a chart's dependencies.
type Resolver struct {
	cfg      *config.Configfile
	workdir  string
	cachedir string
	rescache map[string]*Resolution
}

// Resolution describes a particular dependency as it was resolved.
type Resolution struct {
	// Chartfile is the chart file for this dependency (if found).
	Chartfile *chart.Chartfile
	// Fetched indicates that his chart is in the workspace.
	Fetched bool
	// Found indicates that this chart has been found in workspace or a chart repo.
	Found bool
	// Satisfies indicates that this chart satisfies any known dependencies.
	Satisfies bool
	// If Satisfies is false and Found is true, then this tells how it fails to satisfy.
	SatisfiesErr error
}

// NewResolver creates a new Resolver.
func NewResolver(cfg *config.Configfile, workdir, cachedir string) *Resolver {
	return &Resolver{cfg: cfg, workdir: workdir, cachedir: cachedir}
}

// Resolve starts with a chart and finds all of its dependencies.
func (r *Resolver) Resolve(ch *chart.Chartfile, repodir string) (map[string]*Resolution, error) {
	// We do this here to allow the same resolver to resolve multiple charts.
	r.rescache = map[string]*Resolution{}

	repo := r.cfg.Repos.ByName(repodir)

	// Make it easy to resolve dependencies by setting a From field.
	if ch.From == nil {
		ch.From = &chart.Dependency{Name: ch.Name, Repo: repo}
	}

	key := fmt.Sprintf("%s:%s", repo, ch.Name)
	// Put self in cache (for circular dependencies). It is Found because we have
	// already loaded. It is Satisfied because it is the top-level dep, and it
	// is not marked Fetched. Note that in a circular dependency situation,
	// the Satisfies flag can be set to false if a dependent chart requires
	// a version that this top-level can't satisfy.
	r.rescache[key] = &Resolution{Chartfile: ch, Found: true, Satisfies: true}

	err := r.resolve(ch, ch)
	delete(r.rescache, key)
	return r.rescache, err
}

// repoOrDefault returns the repo that this chart is associated with.
func (r *Resolver) repoOrDefault(ch, parent *chart.Chartfile) (string, error) {
	defRepo := ""
	if ch.From != nil && ch.From.Repo != "" {
		// If this declares its own repo, use it.
		defRepo = ch.From.Repo
	} else if parent.From != nil && parent.From.Repo != "" {
		// IF the parent has a repo, inherit it.
		defRepo = parent.From.Repo
	} else if dr := r.cfg.Repos.Default; dr != "" {
		// Otherwise, get the default repo.
		defRepo = r.cfg.Repos.ByName(dr)
	}
	return canonicalRepo(defRepo)
}

// resolve recursively resolves a chartfile's dependencies.
//
// It does not reset the cache.
func (r *Resolver) resolve(ch, parent *chart.Chartfile) error {

	// A default repo is used to calculate dependencies when a Dependency does
	// not declare a Repo explicitly. In this case, we assume that the
	// Repo is inherited from the parent chart. We check...
	// - From.Repo for a repo, and if it's empty
	// - r.cfg.Repos.Default
	defRepo, err := r.repoOrDefault(ch, parent)
	if err != nil {
		return err
	}

	for _, dep := range ch.Dependencies {
		log.Debug("Testing %s against %s", ch.Name, dep.Name)
		if dep.Repo == "" {
			dep.Repo = defRepo
		}
		var res *Resolution
		// If we have a cached copy of the cache, use that. Otherwise, try to
		// load the chart.
		var ok bool
		if res, ok = r.rescache[dep.Name]; !ok {
			log.Debug("Loading chart %s:%s. This should only happen once.", dep.Repo, dep.Name)
			var err error
			res, err = r.findChart(dep, defRepo)
			if err != nil {
				return err
			}
			r.rescache[dep.Repo+":"+dep.Name] = res
		}

		if res.Found == false {
			continue
		}

		// Check that chart satisfies the current dependency.
		if err := Satisfies(res.Chartfile, dep); err != nil {
			res.Satisfies = false
			res.SatisfiesErr = err
		} else if res.SatisfiesErr == nil {
			// We only set to true if some other run hasn't failed.
			res.Satisfies = true
		}

		// Resolve all of this dependency's dependencies.
		if err := r.resolve(res.Chartfile, ch); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) findChart(d *chart.Dependency, defRepo string) (*Resolution, error) {
	// The order we do this is as follows:
	// - Read all of the charts in the workspace, looking for one whose From.Name/From.Repo matches
	// - If none, read all of the charts in the matching cache.
	// - If none, error

	// Read all charts in the workspace
	ws, err := filepath.Glob(filepath.Join(r.workdir, "*/Chart.yaml"))
	if err != nil {
		return nil, err
	}

	// See if this workspace chart satisfies a dependency.
	for _, path := range ws {
		log.Debug("Loading chart %s", path)
		cf, err := chart.LoadChartfile(path)
		if err != nil {
			log.Warn("Skipping chart %s: %s", path, err)
			continue
		}
		if chartMatches(cf, d) {
			return &Resolution{Chartfile: cf, Fetched: true, Found: true}, nil
		}
	}

	// If we get here, we didn't find a chart in workspace, so we need to find
	// it in the repo that the dependency specifies.
	if d.Repo == "" {
		// Set the repo to the location of origin chart.
		d.Repo = defRepo
	}
	repo, err := r.repoPath(d)
	if err != nil {
		if err == errNoRepoFound {
			return &Resolution{Found: false}, nil
		}
		return nil, err
	}

	// What we're testing here is that Repo and Name match the repository
	// and the chart name.
	path := filepath.Join(repo, "Chart.yaml")
	log.Debug("Loading chart %s", path)
	cf, err := chart.LoadChartfile(path)
	if err != nil {
		return nil, err
	}

	// Simulate a From reference so that the repository information is easy
	// to access.
	cf.From = &chart.Dependency{Name: cf.Name, Repo: d.Repo}
	log.Debug("From: %v (repo: %s)", cf.From, d.Repo)

	// FIXME: Do we need any additional checks?
	// (Version is checked elsewhere)
	return &Resolution{Chartfile: cf, Found: true, Fetched: false}, nil
}

var errNoRepoFound = errors.New("no repo found")

// repoPath gets the path to the repo referenced in d.
func (r *Resolver) repoPath(d *chart.Dependency) (string, error) {
	repo, err := canonicalRepo(d.Repo)
	log.Info("Looking for %s:%s %s", repo, d.Name, d.Version)
	if err != nil {
		log.Warn("Could not calculate canonical repo name %s: %s", d.Repo, err)
		repo = d.Repo
	}
	for _, t := range r.cfg.Repos.Tables {
		trepo, err := canonicalRepo(t.Repo)
		log.Debug("scanning %s", trepo)
		if err != nil {
			log.Warn("Could not calculate canonical repo name %s: %s", t.Repo, err)
			continue
		}
		if repo == trepo {
			return filepath.Join(r.cachedir, t.Name, d.Name), nil
		}
	}
	return "", errNoRepoFound
}

func chartMatches(cf *chart.Chartfile, d *chart.Dependency) bool {
	if cf.From == nil {
		log.Debug("Skipping chart %s: No From.", cf.Name)
		return false
	}

	if cf.From.Name == d.Name {
		if d.Repo == "" {
			log.Debug("Dependency does not have Repo, but %s matches cf.From.Name.", cf.Name)
			return true
		}
		log.Debug("Dependency %s matches name %s", cf.Name, d.Name)
		r := reposMatch(d.Repo, cf.From.Repo)
		if !r {
			log.Debug("Repository %s does not match %s", d.Repo, cf.From.Repo)
		}
		return r
	}
	return false
}

// Satisfies compares the chart to the requirement.
//
// It returns an error if the chart fails to satisfy the requirement. The
// error indicates how the chart failed to satisfy the dependency.
func Satisfies(cf *chart.Chartfile, req *chart.Dependency) error {
	if req.VersionOK(cf.Version) && req.Name == cf.Name {
		return nil
	}

	// I think this is what we want to do.
	if cf.From != nil && req.VersionOK(cf.From.Version) && req.Name == cf.From.Name {
		return nil
	}

	return fmt.Errorf("%s %s does not satisfy %s %s", cf.Name, cf.Version, req.Name, req.Version)
}
