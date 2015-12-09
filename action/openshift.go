package action

import (
	"os/exec"
	"strings"

	"github.com/helm/helm/log"

	"k8s.io/kubernetes/pkg/api/v1"

	// lets force the initialisation of the OAuthClient scheme
	_ "github.com/openshift/origin/pkg/oauth/api/v1"
)

var initialisedOpenShiftFlag = false
var openshiftCluster = false
var openshiftProject = ""

// isOpenShift returns true if the current shell is running against an OpenShift installation of Kubernetes
// or returns false if its a regular Kubernetes platform
func isOpenShift() bool {
	if !initialisedOpenShiftFlag {
		initialisedOpenShiftFlag = true

		cmd := "oc"
		a := []string{"project"}
		b, err := exec.Command(cmd, a...).CombinedOutput()
		if err == nil {
			text := string(b)
			log.Info("got project info: %s", text)
			prefix := "Using project \""
			remaining := strings.TrimPrefix(text, prefix)
			if len(remaining) != len(text) {
				openshiftCluster = true
				idx := strings.Index(remaining, "\"")
				if idx > 0 {
					openshiftProject = remaining[0:idx]
					log.Info("OpenShift is using the project %s", openshiftProject)
				}
			} else {
				log.Debug("Not OpenShift as command: %s %s returned: %s", cmd, strings.Join(a, " "), text)
			}
		} else {
			log.Debug("Not OpenShift as command: %s %s failed with: %s", cmd, strings.Join(a, " "), err)
		}
	}
	return openshiftCluster
}

func createOpenShiftRouteIfRequired(service *v1.Service, ns string, mode string, dryRun bool) error {
	if !isOpenShift() {
		return nil
	}
	t := service.Spec.Type
	if t != "LoadBalancer" {
		return nil
	}
	if ns != "" && ns != openshiftProject {
		ocProject(ns)
		defer ocProject(openshiftProject)
	}
	name := service.Name

	b, err := ocGetRouteYAML(name)
	if err == nil {
		log.Debug("Not trying to expose service %s as we already have a route %s", name, string(b))
		return nil
	}
	cmd := "oc"
	a := []string{"expose", "service", name}
	b, err = exec.Command(cmd, a...).CombinedOutput()
	if err != nil {
		log.Warn("Failed to expose the service %s. Command: `%s %s` failed: %s with %s", name, cmd, strings.Join(a, " "), string(b), err)
	}
	return err
}

func ocProject(ns string) error {
	cmd := "oc"
	a := []string{"project"}
	b, err := exec.Command(cmd, a...).CombinedOutput()
	if err == nil {
		log.Debug("Switched to OpenShift project %s with results %s", ns, string(b))
	} else {
		log.Debug("Failed to set OpenShift project %s due to: %s", ns, err)
	}
	return err
}

func ocGetRouteYAML(name string) ([]byte, error) {
	cmd := "oc"
	a := []string{"get", "route", name, "-oyaml"}
	return exec.Command(cmd, a...).Output()
}
