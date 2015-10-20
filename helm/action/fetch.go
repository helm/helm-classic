package action

import (
	"os"
	"path/filepath"
)

func Fetch(chart, homedir string) {
	src := filepath.Join(homedir, CachePath, chart)
	dest := filepath.Join(homedir, ManifestsPath, chart)

	if fi, err := os.Stat(src); err != nil {
		Die("Could not find %s: %s", chart, err)
	} else if !fi.IsDir() {
		Die("Malformed chart %s: Chart must be in a directory.", chart)
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		Die("Could not create %q: %s", dest, err)
	}

	Info("Fetching %s to %s", src, dest)
	Info("Executing pre-install")
	Info("Templating manifests")
	Info("Copying manifests to %s", dest)
}
