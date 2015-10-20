package action

import (
	"path/filepath"
)

// List lists all of the local charts.
func List(homedir, ns string) {
	if ns == "" {
		ns = DefaultNS
	}
	md := filepath.Join(homedir, ManifestsPath, ns, "*")
	charts, err := filepath.Glob(md)
	if err != nil {
		Warn("Could not find any charts in %q: %s", md, err)
	}
	for _, c := range charts {
		Info("\t%s", c)
	}
}
