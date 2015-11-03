package action

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/deis/helm/chart"
	"github.com/deis/helm/log"
)

func Uninstall(chartName, home, namespace string) {
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
	if err := deleteChart(c, namespace); err != nil {
		log.Die("Failed to completely delete chart: %s", err)
	}
	log.Info("Done")
}

func deleteChart(c *chart.Chart, ns string) error {
	// We delete charts in the ALMOST reverse order that we created them. We
	// start with services to effectively shut down traffic. Then we delete
	// rcs and pods.
	ktype := "service"
	for _, o := range c.Services {
		if err := kubectlDelete(o.Name, ktype, ns); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "rc"
	for _, o := range c.ReplicationControllers {
		if err := kubectlDelete(o.Name, ktype, ns); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "pod"
	for _, o := range c.Pods {
		if err := kubectlDelete(o.Name, ktype, ns); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}

	ktype = "secret"
	for _, o := range c.Secrets {
		if err := kubectlDelete(o.Name, ktype, ns); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "persistentvolume"
	for _, o := range c.PersistentVolumes {
		if err := kubectlDelete(o.Name, ktype, ns); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "namespace"
	for _, o := range c.Namespaces {
		if err := kubectlDelete(o.Name, ktype, ns); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}

	return nil
}

func kubectlDelete(name, ktype, ns string) error {
	log.Debug("Deleting %s (%s)", name, ktype)
	a := []string{"delete", ktype, name}
	if ns != "" {
		a = append([]string{fmt.Sprintf("--namespace=%q", ns)}, a...)
	}

	cmd := exec.Command("kubectl", a...)

	if d, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", string(d), err)
	}
	return nil
}
