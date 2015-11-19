package action

import (
	"io"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/log"
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
func Uninstall(chartName, home, namespace string, force bool) {
	if !chartFetched(chartName, home) {
		log.Info("No chart named %q in your workspace. Nothing to delete.", chartName)
		return
	}

	cd := filepath.Join(home, WorkspaceChartPath, chartName)
	c, err := chart.Load(cd)
	if err != nil {
		log.Die("Failed to load chart: %s", err)
	}
	if err := deleteChart(c, namespace, true); err != nil {
		log.Die("Failed to list charts: %s", err)
	}
	if !force && !promptConfirm("Uninstall the listed objects?") {
		log.Info("Aborted uninstall")
		return
	}

	log.Info("Running `kubectl delete` ...")
	if err := deleteChart(c, namespace, false); err != nil {
		log.Die("Failed to completely delete chart: %s", err)
	}
	log.Info("Done")
}

// promptConfirm prompts a user to confirm (or deny) something.
//
// True is returned iff the prompt is confirmed.
// Errors are reported to the log, and return false.
//
// Valid confirmations:
// 	y, yes, true, t, aye-aye
//
// Valid denials:
//	n, no, f, false
//
// Any other prompt response will return false, and issue a warning to the
// user.
func promptConfirm(msg string) bool {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Err("Could not get terminal: %s", err)
		return false
	}
	defer terminal.Restore(0, oldState)

	f := readerWriter(log.Stdin, log.Stdout)
	t := terminal.NewTerminal(f, msg+" (y/N) ")
	res, err := t.ReadLine()
	if err != nil {
		log.Err("Could not read line: %s", err)
		return false
	}
	res = strings.ToLower(res)
	switch res {
	case "yes", "y", "true", "t", "aye-aye":
		return true
	case "no", "n", "false", "f":
		return false
	}
	log.Warn("Did not understand answer %q, assuming No", res)
	return false
}

func readerWriter(reader io.Reader, writer io.Writer) *rw {
	return &rw{r: reader, w: writer}
}

// rw is a trivial io.ReadWriter that does not buffer.
type rw struct {
	r io.Reader
	w io.Writer
}

func (x *rw) Read(b []byte) (int, error) {
	return x.r.Read(b)
}
func (x *rw) Write(b []byte) (int, error) {
	return x.w.Write(b)
}

func deleteChart(c *chart.Chart, ns string, dry bool) error {
	// We delete charts in the ALMOST reverse order that we created them. We
	// start with services to effectively shut down traffic.
	ktype := "service"
	for _, o := range c.Services {
		if dry {
			log.Msg("%s/%s", ktype, o.Name)
		} else if _, err := Kubectl.Delete(o.Name, ktype, ns, false); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "pod"
	for _, o := range c.Pods {
		if dry {
			log.Msg("%s/%s", ktype, o.Name)
		} else if _, err := Kubectl.Delete(o.Name, ktype, ns, false); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "rc"
	for _, o := range c.ReplicationControllers {
		if dry {
			log.Msg("%s/%s", ktype, o.Name)
		} else if _, err := Kubectl.Delete(o.Name, ktype, ns, false); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "secret"
	for _, o := range c.Secrets {
		if dry {
			log.Msg("%s/%s", ktype, o.Name)
		} else if _, err := Kubectl.Delete(o.Name, ktype, ns, false); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "persistentvolume"
	for _, o := range c.PersistentVolumes {
		if dry {
			log.Msg("%s/%s", ktype, o.Name)
		} else if _, err := Kubectl.Delete(o.Name, ktype, ns, false); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}
	ktype = "namespace"
	for _, o := range c.Namespaces {
		if dry {
			log.Msg("%s/%s", ktype, o.Name)
		} else if _, err := Kubectl.Delete(o.Name, ktype, ns, false); err != nil {
			log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
		}
	}

	return nil
}
