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

	src := filepath.Join(homedir, ChartPath, chart)
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

	walker := func(fname string, fi os.FileInfo, e error) error {
		if e != nil {
			Warn("Encounter error walking %q: %s", fname, e)
			return nil
		}

		Info("Copying %s", fname)
		rf, err := filepath.Rel(src, fname)
		if err != nil {
			Warn("Could not find relative path: %s", err)
			return nil
		}
		df := filepath.Join(dst, rf)

		// Handle directories by creating mirrors.
		if fi.IsDir() {
			if err := os.MkdirAll(df, fi.Mode()); err != nil {
				Warn("Could not create %q: %s", df, err)
			}
			return nil
		}

		// Otherwise, copy files.
		in, err := os.Open(fname)
		if err != nil {
			Warn("Skipping file %s: %s", fname, err)
			return nil
		}
		out, err := os.Create(df)
		if err != nil {
			in.Close()
			Warn("Skipping file copy %s: %s", fname, err)
			return nil
		}
		if _, err = io.Copy(out, in); err != nil {
			Warn("Copy from %s to %s failed: %s", fname, df, err)
		}

		if err := out.Close(); err != nil {
			Warn("Failed to close %q: %s", df, err)
		}
		if err := in.Close(); err != nil {
			Warn("Failed to close reader %q: %s", fname, err)
		}

		return nil
	}
	return filepath.Walk(src, walker)
}
