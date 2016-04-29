package action

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/manifest"
	helm "github.com/helm/helm-classic/util"
)

// kubeGetter wraps the kubectl command, override in tests
type kubeGetter func(string) string

var kubeGet kubeGetter = func(m string) string {
	log.Debug("Getting manifests from %s", m)

	a := []string{"get", "-f", m}
	out, _ := exec.Command("kubectl", a...).CombinedOutput()
	return string(out)
}

// Remove removes a chart from the workdir.
//
// - chart is the source
// - homedir is the home directory for the user
// - force will remove installed charts from workspace
func Remove(chart, homedir string, force bool) {
	chartPath := helm.WorkspaceChartDirectory(homedir, chart)
	if _, err := os.Stat(chartPath); err != nil {
		log.Err("Chart not found. %s", err)
		return
	}

	if !force {
		var connectionFailure bool

		// check if any chart manifests are installed
		installed, err := checkManifests(chartPath)
		if err != nil {
			if strings.Contains(err.Error(), "unable to connect") {
				connectionFailure = true
			} else {
				log.Die(err.Error())
			}
		}

		if connectionFailure {
			log.Err("Could not determine if %s is installed.  To remove the chart --force flag must be set.", chart)
			return
		} else if len(installed) > 0 {
			log.Err("Found %d installed manifests for %s.  To remove a chart that has been installed the --force flag must be set.", len(installed), chart)
			return
		}
	}

	// remove local chart files
	if err := os.RemoveAll(chartPath); err != nil {
		log.Die("Could not remove chart. %s", err)
	}

	log.Info("All clear! You have successfully removed %s from your workspace.", chart)
}

// checkManifests gets any installed manifests within a chart
func checkManifests(chartPath string) ([]string, error) {
	var foundManifests []string

	manifests, err := manifest.Files(chartPath)
	if err != nil {
		return nil, err
	}

	for _, m := range manifests {
		out := kubeGet(m)

		if strings.Contains(out, "unable to connect") {
			return nil, fmt.Errorf(out)
		}
		if !strings.Contains(out, "not found") {
			foundManifests = append(foundManifests, m)
		}
	}

	log.Debug("Found %d installed manifests", len(foundManifests))

	return foundManifests, nil
}
