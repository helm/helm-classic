package action

import (
	"io"
	"os"
	"path/filepath"
)

func Fetch(chart, lname, homedir string) {

	if lname == "" {
		lname = chart
	}

	src := filepath.Join(homedir, CachePath, chart)
	dest := filepath.Join(homedir, ManifestsPath, lname)

	if fi, err := os.Stat(src); err != nil {
		Die("Could not find %s: %s", chart, err)
	} else if !fi.IsDir() {
		Die("Malformed chart %s: Chart must be in a directory.", chart)
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		Die("Could not create %q: %s", dest, err)
	}

	Info("Fetching %s to %s", src, dest)
	if err := copyDir(src, dest); err != nil {
	}
	Info("Executing pre-install")
	Info("Templating manifests")
	Info("Copying manifests to %s", dest)
}

// Copy a directory and its subdirectories.
func copyDir(src, dst string) error {
	files, err := filepath.Glob(filepath.Join(src, "*"))
	if err != nil {
		return err
	}

	for _, fname := range files {
		Info("Copying %s", fname)
		rf, err := filepath.Rel(src, fname)
		if err != nil {
			Warn("Could not find relative path: %s", err)
			continue
		}

		df := filepath.Join(dst, rf)
		in, err := os.Open(src)
		if err != nil {
			Warn("Skipping file %s: %s", src, err)
			continue
		}
		out, err := os.Create(df)
		if err != nil {
			in.Close()
			Warn("Skipping file copy %s: %s", src, err)
			continue
		}
		if _, err = io.Copy(out, in); err != nil {
			Warn("Copy from %s to %s failed: %s", src, df, err)
		}

		out.Close()
		in.Close()
	}
	return nil
}
