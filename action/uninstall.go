package action

import (
	"io"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/helm/helm-classic/chart"
	"github.com/helm/helm-classic/kubectl"
	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/manifest"
	helm "github.com/helm/helm-classic/util"
)

// Uninstall removes a chart from Kubernetes.
//
// Manifests are removed from Kubernetes in the order specified by
// chart.UninstallOrder. Any unknown types are removed before that sequence
// is run.
func Uninstall(chartName, home, namespace string, force bool, client kubectl.Runner) {
	// This is a stop-gap until kubectl respects namespaces in manifests.
	if namespace == "" {
		log.Die("This command requires a namespace. Did you mean '-n default'?")
	}
	if !chartFetched(chartName, home) {
		log.Info("No chart named %q in your workspace. Nothing to delete.", chartName)
		return
	}

	cd := helm.WorkspaceChartDirectory(home, chartName)
	c, err := chart.Load(cd)
	if err != nil {
		log.Die("Failed to load chart: %s", err)
	}
	if err := deleteChart(c, namespace, true, client); err != nil {
		log.Die("Failed to list charts: %s", err)
	}
	if !force && !promptConfirm("Uninstall the listed objects?") {
		log.Info("Aborted uninstall")
		return
	}

	CheckKubePrereqs()

	log.Info("Running `kubectl delete` ...")
	if err := deleteChart(c, namespace, false, client); err != nil {
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

// deleteChart deletes all of the Kubernetes manifests associated with this chart.
func deleteChart(c *chart.Chart, ns string, dry bool, client kubectl.Runner) error {
	// Unknown kinds get uninstalled first because we know that core kinds
	// do not depend on them.
	for _, kind := range c.UnknownKinds(UninstallOrder) {
		uninstallKind(c.Kind[kind], ns, kind, dry, client)
	}

	// Uninstall all of the known kinds in a particular order.
	for _, kind := range UninstallOrder {
		uninstallKind(c.Kind[kind], ns, kind, dry, client)
	}

	return nil
}

func uninstallKind(kind []*manifest.Manifest, ns, ktype string, dry bool, client kubectl.Runner) {
	for _, o := range kind {
		if dry {
			log.Msg("%s/%s", ktype, o.Name)
		} else {
			// If it's a keeper manifest, skip uninstall.
			if data, err := o.VersionedObject.JSON(); err == nil {
				if manifest.IsKeeper(data) {
					log.Warn("Not uninstalling %s %s because of \"helm-keep\" annotation.\n"+
						"---> Use kubectl to uninstall keeper manifests.\n", ktype, o.Name)
					continue
				}
			}
			out, err := client.Delete(o.Name, ktype, ns)
			if err != nil {
				log.Warn("Could not delete %s %s (Skipping): %s", ktype, o.Name, err)
			}
			log.Info(string(out))
		}
	}
}
