package action

import (
	"path/filepath"

	"github.com/deis/helm/helm/model"
)

// List lists all of the local charts.
func List(homedir, ns string) {
	if ns == "" {
		ns = DefaultNS
	}

	// List all namespaces
	if ns == "*" {
		md := filepath.Join(homedir, ManifestsPath, ns)
		nss, err := filepath.Glob(md)
		if err != nil {
			Warn("Could not find any namespaces in %q: %s", md, err)
		}
		for _, n := range nss {
			dir := filepath.Base(n)
			Info("%s:", dir)
			listNS(homedir, dir)
		}
		return
	}
	listNS(homedir, ns)
}

func listNS(homedir, ns string) {
	md := filepath.Join(homedir, ManifestsPath, ns, "*")
	charts, err := filepath.Glob(md)
	if err != nil {
		Warn("Could not find any charts in %q: %s", md, err)
	}
	for _, c := range charts {
		cname := filepath.Base(c)
		if ch, err := model.Load(filepath.Join(c, "Chart.yaml")); err == nil {
			Info("\t%s (%s %s) - %s", cname, ch.Name, ch.Version, ch.Description)
			continue
		}
		Info("\t%s (unknown)", cname)
	}
}
