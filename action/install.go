package action

import (
	"bytes"
	"os"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/codec"
	"github.com/helm/helm/dependency"
	"github.com/helm/helm/kubectl"
	"github.com/helm/helm/log"
	helm "github.com/helm/helm/util"
)

// InstallOrder defines the order in which manifests should be installed, by Kind.
//
// Anything not on the list will be installed after the last listed item, in
// an indeterminate order.
var InstallOrder = []string{"Namespace", "Secret", "PersistentVolume", "ServiceAccount", "Service", "Pod", "ReplicationController", "DaemonSet", "Ingress", "Job"}

// UninstallOrder defines the order in which manifests are uninstalled.
//
// Unknown manifest types (those not explicitly referenced in this list) will
// be uninstalled before any of these, since we know that none of the core
// types depend on non-core types.
var UninstallOrder = []string{"Service", "Pod", "ReplicationController", "DaemonSet", "Secret", "PersistentVolume", "ServiceAccount", "Ingress", "Job", "Namespace"}

// Install loads a chart into Kubernetes.
//
// If the chart is not found in the workspace, it is fetched and then installed.
//
// During install, manifests are sent to Kubernetes in the ordered specified by InstallOrder.
func Install(chartName, home, namespace string, force bool, generate bool, exclude []string, client kubectl.Runner) {
	ochart := chartName
	r := mustConfig(home).Repos
	table, chartName := r.RepoChart(chartName)

	if !chartFetched(chartName, home) {
		log.Info("No chart named %q in your workspace. Fetching now.", ochart)
		fetch(chartName, chartName, home, table)
	}

	cd := helm.WorkspaceChartDirectory(home, chartName)
	c, err := chart.Load(cd)
	if err != nil {
		log.Die("Failed to load chart: %s", err)
	}

	// Give user the option to bale if dependencies are not satisfied.
	nope, err := dependency.Resolve(c.Chartfile, helm.WorkspaceChartDirectory(home))

	if err != nil {
		log.Warn("Failed to check dependencies: %s", err)
		if !force {
			log.Die("Re-run with --force to install anyway.")
		}
	} else if len(nope) > 0 {
		log.Warn("Unsatisfied dependencies:")
		for _, d := range nope {
			log.Msg("\t%s %s", d.Name, d.Version)
		}
		if !force {
			log.Die("Stopping install. Re-run with --force to install anyway.")
		}
	}

	// Run the generator if -g is set.
	if generate {
		Generate(chartName, home, exclude, force)
	}

	CheckKubePrereqs()

	log.Info("Running `kubectl create -f` ...")
	if err := uploadManifests(c, namespace, client); err != nil {
		log.Die("Failed to upload manifests: %s", err)
	}
	log.Info("Done")

	PrintREADME(chartName, home)
}

// uploadManifests sends manifests to Kubectl in a particular order.
func uploadManifests(c *chart.Chart, namespace string, client kubectl.Runner) error {

	// Install known kinds in a predictable order.
	for _, k := range InstallOrder {
		for _, m := range c.Kind[k] {
			o := m.VersionedObject
			o.AddAnnotations(map[string]string{
				chart.AnnFile:         m.Source,
				chart.AnnChartVersion: c.Chartfile.Version,
				chart.AnnChartDesc:    c.Chartfile.Description,
				chart.AnnChartName:    c.Chartfile.Name,
			})
			var data []byte
			var err error
			if data, err = o.JSON(); err != nil {
				return err
			}
			log.Debug("File: %s", string(data))
			out, err := client.Create(data, namespace)
			log.Msg(string(out))
			if err != nil {
				return err
			}
		}
	}

	// Install unknown kinds afterward. Order here is not predictable.
	for _, k := range c.UnknownKinds(InstallOrder) {
		for _, o := range c.Kind[k] {
			o.VersionedObject.AddAnnotations(map[string]string{chart.AnnFile: o.Source})
			out, err := marshalAndCreate(o.VersionedObject, namespace, client)
			log.Msg(string(out))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func marshalAndCreate(o interface{}, ns string, client kubectl.Runner) ([]byte, error) {
	var b bytes.Buffer
	if err := codec.JSON.Encode(&b).One(o); err != nil {
		return nil, err
	}
	data := b.Bytes()
	log.Debug("File: %s", string(data))
	return client.Create(data, ns)
}

// Check by chart directory name whether a chart is fetched into the workspace.
//
// This does NOT check the Chart.yaml file.
func chartFetched(chartName, home string) bool {
	p := helm.WorkspaceChartDirectory(home, chartName, Chartfile)
	log.Debug("Looking for %q", p)
	if fi, err := os.Stat(p); err != nil || fi.IsDir() {
		log.Debug("No chart: %s", err)
		return false
	}
	return true
}
