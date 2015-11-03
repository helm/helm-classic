package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/deis/helm/dependency"
	"github.com/deis/helm/manifest"
	"github.com/deis/helm/model"
)

import (
	"github.com/deis/helm/log"
)

// Install loads a chart into Kubernetes.
//
// If the chart is not found in the workspace, it is fetched and then installed.
//
// During install, manifests are sent to Kubernetes in the following order:
//
//	- Namespaces
// 	- Secrets
// 	- Volumes
// 	- Services
// 	- Pods
// 	- ReplicationControllers
func Install(chart, home, namespace string, force bool) {
	if !chartInstalled(chart, home) {
		log.Info("No installed chart named %q. Installing now.", chart)
		fetch(chart, chart, home)
	}

	cd := filepath.Join(home, WorkspaceChartPath, chart)
	c, err := model.Load(cd)
	if err != nil {
		log.Die("Failed to load chart: %s", err)
	}

	// Give user the option to bale if dependencies are not satisfied.
	nope, err := dependency.Resolve(c.Chartfile, filepath.Join(home, WorkspaceChartPath))
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

	if err := uploadManifests(c, namespace); err != nil {
		log.Die("Failed to upload manifests: %s", err)
	}
	PrintREADME(chart, home)
}

func AltInstall(chart, cachedir, home, namespace string, force bool) {
	// Make sure there is a chart in the cachedir.
	if _, err := os.Stat(filepath.Join(cachedir, "Chart.yaml")); err != nil {
		log.Die("Expected a Chart.yaml in %s: %s", cachedir, err)
	}
	// Make sure there is a manifests dir.
	if fi, err := os.Stat(filepath.Join(cachedir, "manifests")); err != nil {
		log.Die("Expected 'manifests/' in %s: %s", cachedir, err)
	} else if !fi.IsDir() {
		log.Die("Expected 'manifests/' to be a directory in %s: %s", cachedir, err)
	}

	// Copy the source chart to the workspace. We ruthlessly overwrite in
	// this case.
	dest := filepath.Join(home, WorkspaceChartPath, chart)
	if err := copyDir(cachedir, dest); err != nil {
		log.Die("Failed to copy %s to %s: %s", cachedir, dest, err)
	}

	// Load the chart.
	c, err := model.Load(dest)
	if err != nil {
		log.Die("Failed to load chart: %s", err)
	}

	// Give user the option to bale if dependencies are not satisfied.
	nope, err := dependency.Resolve(c.Chartfile, filepath.Join(home, WorkspaceChartPath))
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

	if err := uploadManifests(c, namespace); err != nil {
		log.Die("Failed to upload manifests: %s", err)
	}
}

// uploadManifests sends manifests to Kubectl in a particular order.
func uploadManifests(c *model.Chart, namespace string) error {
	// The ordering is significant.
	// TODO: Right now, we force version v1. We could probably make this more
	// flexible if there is a use case.
	for _, o := range c.Namespaces {
		b, err := manifest.MarshalJSON(o, "v1")
		if err != nil {
			return err
		}
		if err := kubectlCreate(b, namespace); err != nil {
			return err
		}
	}
	for _, o := range c.Secrets {
		b, err := manifest.MarshalJSON(o, "v1")
		if err != nil {
			return err
		}
		if err := kubectlCreate(b, namespace); err != nil {
			return err
		}
	}
	for _, o := range c.PersistentVolumes {
		b, err := manifest.MarshalJSON(o, "v1")
		if err != nil {
			return err
		}
		if err := kubectlCreate(b, namespace); err != nil {
			return err
		}
	}
	for _, o := range c.Services {
		b, err := manifest.MarshalJSON(o, "v1")
		if err != nil {
			return err
		}
		if err := kubectlCreate(b, namespace); err != nil {
			return err
		}
	}
	for _, o := range c.Pods {
		b, err := manifest.MarshalJSON(o, "v1")
		if err != nil {
			return err
		}
		if err := kubectlCreate(b, namespace); err != nil {
			return err
		}
	}
	for _, o := range c.ReplicationControllers {
		b, err := manifest.MarshalJSON(o, "v1")
		if err != nil {
			return err
		}
		if err := kubectlCreate(b, namespace); err != nil {
			return err
		}
	}
	return nil
}

// Check by chart directory name whether a chart is installed.
//
// This does NOT check the Chart.yaml file.
func chartInstalled(chart, home string) bool {
	p := filepath.Join(home, WorkspaceChartPath, chart, "Chart.yaml")
	log.Debug("Looking for %q", p)
	if fi, err := os.Stat(p); err != nil || fi.IsDir() {
		log.Debug("No chart: %s", err)
		return false
	}
	return true
}

// kubectlCreate calls `kubectl create` and sends the data via Stdin.
func kubectlCreate(data []byte, ns string) error {
	a := []string{"create", "-f", "-"}

	if ns != "" {
		a = append([]string{"--namespace=" + ns}, a...)
	}

	c := exec.Command("kubectl", a...)
	in, err := c.StdinPipe()
	if err != nil {
		return err
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		return err
	}

	log.Info("File: %s", string(data))
	in.Write(data)
	in.Close()

	return c.Wait()
}
