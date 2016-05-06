package action

import (
	"os"
	"path/filepath"

	"github.com/helm/helm-classic/chart"
	"github.com/helm/helm-classic/dependency"
	"github.com/helm/helm-classic/log"
	helm "github.com/helm/helm-classic/util"
)

// Fetch gets a chart from the source repo and copies to the workdir.
//
// - chartName is the source
// - lname is the local name for that chart (chart-name); if blank, it is set to the chart.
// - homedir is the home directory for the user
func Fetch(chartName, lname, homedir string) {

	r := mustConfig(homedir).Repos
	repository, chartName := r.RepoChart(chartName)

	if lname == "" {
		lname = chartName
	}

	fetch(chartName, lname, homedir, repository)

	chartFilePath := helm.WorkspaceChartDirectory(homedir, lname, Chartfile)
	cfile, err := chart.LoadChartfile(chartFilePath)
	if err != nil {
		log.Die("Source is not a valid chart. Missing Chart.yaml: %s", err)
	}

	deps, err := dependency.Resolve(cfile, helm.WorkspaceChartDirectory(homedir))
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

	log.Info("Fetched chart into workspace %s", helm.WorkspaceChartDirectory(homedir, lname))
	log.Info("Done")
}

func fetch(chartName, lname, homedir, chartpath string) {
	src := helm.CacheDirectory(homedir, chartpath, chartName)
	dest := helm.WorkspaceChartDirectory(homedir, lname)

	fi, err := os.Stat(src)
	if err != nil {
		log.Warn("Oops. Looks like there was an issue finding the chart, %s, in %s. Running `helmc update` to ensure you have the latest version of all Charts from Github...", lname, src)
		Update(homedir)
		fi, err = os.Stat(src)
		if err != nil {
			log.Die("Chart %s not found in %s", lname, src)
		}
		log.Info("Good news! Looks like that did the trick. Onwards and upwards!")
	}

	if !fi.IsDir() {
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

func updateChartfile(src, dest, lname string) error {
	sc, err := chart.LoadChartfile(filepath.Join(src, Chartfile))
	if err != nil {
		return err
	}

	dc, err := chart.LoadChartfile(filepath.Join(dest, Chartfile))
	if err != nil {
		return err
	}

	dc.Name = lname
	dc.From = &chart.Dependency{
		Name:    sc.Name,
		Version: sc.Version,
		Repo:    chart.RepoName(src),
	}

	return dc.Save(filepath.Join(dest, Chartfile))
}
