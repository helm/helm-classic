package action

import (
	"io"
	"os"
	"path/filepath"

	"github.com/deis/helm/chart"
	"github.com/deis/helm/dependency"
	"github.com/deis/helm/log"
	"github.com/deis/helm/repo"
)

// Fetch gets a chart from the source repo and copies to the workdir.
//
// - chartName is the source
// - lname is the local name for that chart (chart-name); if blank, it is set to the chart.
// - homedir is the home directory for the user
// - ns is the namespace for this package. If blank, it is set to the DefaultNS.
func Fetch(chartName, lname, homedir string) {

	r := mustRepofile(homedir)
	repository, chartName := r.RepoChart(chartName)

	if lname == "" {
		lname = chartName
	}

	fetch(chartName, lname, homedir, repository)

	chartFilePath := filepath.Join(homedir, WorkspaceChartPath, lname, "Chart.yaml")
	cfile, err := chart.LoadChartfile(chartFilePath)
	if err != nil {
		log.Die("Source is not a valid chart. Missing Chart.yaml: %s", err)
	}

	deps, err := dependency.Resolve(cfile, filepath.Join(homedir, WorkspaceChartPath))
	if err != nil {
		log.Warn("Could not check dependencies: %s", err)
		return
	}

	if len(deps) > 0 {
		log.Warn("Unsatisfied dependencies:")
		for _, d := range deps {
			log.Msg("\t%s %s", d.Name, d.Version)
		}
	}

	log.Info("Fetched chart into workspace %s", filepath.Join(homedir, WorkspaceChartPath, lname))
	log.Info("Done")
}

// mustRepofile attempts to laod the Repofile or dies trying.
func mustRepofile(homedir string) *repo.Repofile {
	rpath := filepath.Join(homedir, CachePath, Repofile)
	r, err := repo.LoadRepofile(rpath)
	if err != nil {
		log.Die("Could not load %s: %s", rpath, err)
	}
	return r
}

func fetch(chartName, lname, homedir, chartpath string) {
	src := filepath.Join(homedir, CachePath, chartpath, chartName)
	dest := filepath.Join(homedir, WorkspaceChartPath, lname)

	if fi, err := os.Stat(src); err != nil {
	} else if !fi.IsDir() {
		log.Die("Malformed chart %s: Chart must be in a directory.", chartName)
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		log.Die("Could not create %q: %s", dest, err)
	}

	log.Debug("Fetching %s to %s", src, dest)
	if err := copyDir(src, dest); err != nil {
		log.Die("Failed copying %s to %s", src, dest)
	}
}

// Copy a directory and its subdirectories.
func copyDir(src, dst string) error {

	var failure error

	walker := func(fname string, fi os.FileInfo, e error) error {
		if e != nil {
			log.Warn("Encounter error walking %q: %s", fname, e)
			failure = e
			return nil
		}

		log.Debug("Copying %s", fname)
		rf, err := filepath.Rel(src, fname)
		if err != nil {
			log.Warn("Could not find relative path: %s", err)
			return nil
		}
		df := filepath.Join(dst, rf)

		// Handle directories by creating mirrors.
		if fi.IsDir() {
			if err := os.MkdirAll(df, fi.Mode()); err != nil {
				log.Warn("Could not create %q: %s", df, err)
				failure = err
			}
			return nil
		}

		// Otherwise, copy files.
		in, err := os.Open(fname)
		if err != nil {
			log.Warn("Skipping file %s: %s", fname, err)
			return nil
		}
		out, err := os.Create(df)
		if err != nil {
			in.Close()
			log.Warn("Skipping file copy %s: %s", fname, err)
			return nil
		}
		if _, err = io.Copy(out, in); err != nil {
			log.Warn("Copy from %s to %s failed: %s", fname, df, err)
		}

		if err := out.Close(); err != nil {
			log.Warn("Failed to close %q: %s", df, err)
		}
		if err := in.Close(); err != nil {
			log.Warn("Failed to close reader %q: %s", fname, err)
		}

		return nil
	}
	filepath.Walk(src, walker)
	return failure
}
