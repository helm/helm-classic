package action

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

import (
	"github.com/deis/helm/helm/log"
)

func Install(chart, home, namespace string) {
	Fetch(chart, chart, home)
	log.Info("kubectl --namespace=%q create -f %s.yaml", namespace, chart)
	if !chartInstalled(chart, home) {
		log.Info("No installed chart named %q. Installing now.", chart)
		Fetch(chart, chart, home)
	}

	d := filepath.Join(home, WorkspaceChartPath, chart, "manifests")
	log.Debug("Looking for manifests in %q", d)
	files, err := manifestFiles(d)
	if err != nil {
		log.Die("No manifests to install: %s", err)
	}

	for _, f := range files {
		if err := kubectlCreate(f, namespace); err != nil {
			log.Warn("Failed to install manifest %q: %s", f, err)
		}
	}
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

func kubectlCreate(chart, ns string) error {
	a := []string{"create", "-f", chart}
	if ns != "" {
		a = append([]string{fmt.Sprintf("--namespace=%q", ns)}, a...)
	}
	//Info("kubectl --namespace=%q create -f %s", ns, chart)
	c := exec.Command("kubectl", a...)

	// The default error message is not helpful, so we grab the output
	// and prepend it to the error message.
	if o, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", string(o), err)
	}

	return nil
}

func manifestFiles(dir string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, "*.yaml"))
}
