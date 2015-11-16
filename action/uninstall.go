package action

import (
	"path/filepath"

	"github.com/deis/helm/chart"
	"github.com/deis/helm/kubectl"
	"github.com/deis/helm/log"
)

// Uninstall removes a chart from Kubernetes.
//
// Manifests are removed from Kubernetes in the following order:
//
// 	- Services (to shut down traffic)
// 	- Pods (which can be part of RCs)
// 	- ReplicationControllers
// 	- Volumes
// 	- Secrets
//	- Namespaces
func Uninstall(chartName, home, namespace string, dryRun bool) {
	if !chartFetched(chartName, home) {
		log.Info("No chart named %q in your workspace. Nothing to delete.", chartName)
		return
	}

	cd := filepath.Join(home, WorkspaceChartPath, chartName)
	c, err := chart.Load(cd)
	if err != nil {
		log.Die("Failed to load chart: %s", err)
	}

	log.Info("Running `kubectl delete` ...")
	if err := deleteChart(c, namespace, dryRun); err != nil {
		log.Die("Failed to completely delete chart: %s", err)
	}
	log.Info("Done")
}

func deleteChart(c *chart.Chart, ns string, dryRun bool) error {
	// We delete charts in the ALMOST reverse order that we created them. We
	// start with services to effectively shut down traffic.
	ktype := "service"
	for _, o := range c.Services {
		if err := kubectl.Delete(o.Name, ktype, ns, dryRun); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "pod"
	for _, o := range c.Pods {
		if err := kubectl.Delete(o.Name, ktype, ns, dryRun); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "rc"
	for _, o := range c.ReplicationControllers {
		if err := kubectl.Delete(o.Name, ktype, ns, dryRun); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "secret"
	for _, o := range c.Secrets {
		if err := kubectl.Delete(o.Name, ktype, ns, dryRun); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "persistentvolume"
	for _, o := range c.PersistentVolumes {
		if err := kubectl.Delete(o.Name, ktype, ns, dryRun); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "namespace"
	for _, o := range c.Namespaces {
		if err := kubectl.Delete(o.Name, ktype, ns, dryRun); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}

	return nil
}
