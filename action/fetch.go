package action

import (
	"os"
	"path/filepath"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/config"
	"github.com/helm/helm/dependency"
	"github.com/helm/helm/log"
	helm "github.com/helm/helm/util"
)

// Fetch gets a chart from the source repo and copies to the workdir.
//
// - chartName is the source
// - lname is the local name for that chart (chart-name); if blank, it is set to the chart.
// - homedir is the home directory for the user
func Fetch(chartName, lname, homedir string, force bool) {

	cfg := mustConfig(homedir)
	repository, chartName := cfg.Repos.RepoChart(chartName)

	if lname == "" {
		lname = chartName
	}

	chartFilePath := helm.CacheDirectory(homedir, repository, chartName, "Chart.yaml")
	log.Debug("Loading %s", chartFilePath)
	cfile, err := chart.LoadChartfile(chartFilePath)
	if err != nil {
		log.Die("Source is not a valid chart. Missing Chart.yaml: %s", err)
	}

	toFetch := getFetchDependencies(cfg, cfile, homedir, repository, force)
	for _, d := range toFetch {
		log.Info("⇓ Fetching copy of %s %s from %s", d.Chartfile.Name, d.Chartfile.Version, d.Chartfile.From.Repo)
		rn := cfg.Repos.ByRepo(d.Chartfile.From.Repo)
		dn := rn + "/" + d.Chartfile.Name
		ln := rn + "-" + d.Chartfile.Name
		fetch(dn, ln, homedir, rn)
	}
	fetch(chartName, lname, homedir, repository)

	log.Info("Fetched chart into workspace %s", helm.WorkspaceChartDirectory(homedir, lname))
	log.Info("Done")
}

func fetch(chartName, lname, homedir, chartpath string) {
	src := helm.CacheDirectory(homedir, chartpath, chartName)
	dest := helm.WorkspaceChartDirectory(homedir, lname)

	if fi, err := os.Stat(src); err != nil {
		log.Die("Chart %s not found in %s", lname, src)
	} else if !fi.IsDir() {
		log.Die("Malformed chart %s: Chart must be in a directory.", chartName)
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		log.Die("Could not create %q: %s", dest, err)
	}

	log.Debug("Fetching %s to %s", src, dest)
	if err := helm.CopyDir(src, dest); err != nil {
		log.Die("Failed copying %s to %s", src, dest)
	}

	if err := updateChartfile(src, dest, lname); err != nil {
		log.Die("Failed to update Chart.yaml: %s", err)
	}
}

func getFetchDependencies(cfg *config.Configfile, cfile *chart.Chartfile, homedir, repository string, force bool) []*dependency.Resolution {
	workdir := helm.WorkspaceChartDirectory(homedir)
	cachedir := helm.CacheDirectory(homedir)

	cfile.From = &chart.Dependency{
		Repo: cfg.Repos.ByName(repository),
		Name: cfile.Name,
	}

	resolver := dependency.NewResolver(cfg, workdir, cachedir)
	deps, err := resolver.Resolve(cfile, repository)
	if err != nil {
		log.Die("Could not check dependencies: %s", err)
	}

	toFetch := []*dependency.Resolution{}
	already := []*dependency.Resolution{}

	failed := 0
	if len(deps) > 0 {
		for n, d := range deps {
			if d.Found == false {
				failed++
				log.Err("Dependency cannot be found: %s", n)
				continue
			}
			if d.Satisfies == false {
				failed++
				log.Err("Insufficient: %s %s (%s): %s (Workspace chart: %b)", d.Chartfile.Name, d.Chartfile.Version, n, d.SatisfiesErr, d.Fetched)
				if d.Fetched == true {
					log.Info("Re-fetching %s from the chart repo may resolve this.", d.Chartfile.Name)
				}
				continue
			}
			if d.Fetched {
				already = append(already, d)
				continue
			}
			toFetch = append(toFetch, d)
		}
	}
	// If the talied errors are greater than 0, we stop here.
	if failed > 0 {
		if force {
			log.Warn("Ignoring %d errors.", failed)
		} else {
			log.Die("Cannot continue. %d errors.", failed)
		}
	}

	for _, d := range already {
		log.Info("✔︎ Using fetched version of %s %s", d.Chartfile.Name, d.Chartfile.Version)
	}

	return toFetch
}

func updateChartfile(src, dest, lname string) error {
	sc, err := chart.LoadChartfile(filepath.Join(src, "Chart.yaml"))
	if err != nil {
		return err
	}

	dc, err := chart.LoadChartfile(filepath.Join(dest, "Chart.yaml"))
	if err != nil {
		return err
	}

	dc.Name = lname
	dc.From = &chart.Dependency{
		Name:    sc.Name,
		Version: sc.Version,
		Repo:    chart.RepoName(src),
	}

	return dc.Save(filepath.Join(dest, "Chart.yaml"))
}
