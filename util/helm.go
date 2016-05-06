package util

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/helm/helm-classic/log"
)

// Configfile is the file containing Helm Classic's YAML configuration data.
const Configfile = "config.yaml"

// DefaultConfigfile is the default Helm Classic configuration.
const DefaultConfigfile = `repos:
  default: charts
  tables:
    - name: charts
      repo: https://github.com/helm/charts
workspace:
`

// EnsureHome ensures that a HELMC_HOME exists.
func EnsureHome(home string) {

	must := []string{home, CacheDirectory(home), filepath.Join(home, workspacePath)}

	for _, p := range must {
		if fi, err := os.Stat(p); err != nil {
			log.Debug("Creating %s", p)
			if err := os.MkdirAll(p, 0755); err != nil {
				log.Die("Could not create %q: %s", p, err)
			}
		} else if !fi.IsDir() {
			log.Die("%s must be a directory.", home)
		}
	}

	refi := filepath.Join(home, Configfile)
	if _, err := os.Stat(refi); err != nil {
		log.Info("Creating %s", refi)
		// Attempt to create a Repos.yaml
		if err := ioutil.WriteFile(refi, []byte(DefaultConfigfile), 0755); err != nil {
			log.Die("Could not create %s: %s", refi, err)
		}
	}

	if err := os.Chdir(home); err != nil {
		log.Die("Could not change to directory %q: %s", home, err)
	}
}

// CopyDir copy a directory and its subdirectories.
func CopyDir(src, dst string) error {

	var failure error

	walker := func(fname string, fi os.FileInfo, e error) error {
		if e != nil {
			log.Warn("Encounter error walking %q: %s", fname, e)
			failure = e
			return nil
		}

		rf, err := filepath.Rel(src, fname)
		if err != nil {
			log.Warn("Could not find relative path: %s", err)
			return nil
		}
		df := filepath.Join(dst, rf)

		// Handle directories by creating mirrors.
		if fi.IsDir() {
			if err := os.MkdirAll(df, fi.Mode()); err != nil {
				log.Warn("Could not create %q: %s", df, err)
				failure = err
			}
			return nil
		}

		// Otherwise, copy files.
		in, err := os.Open(fname)
		if err != nil {
			log.Warn("Skipping file %s: %s", fname, err)
			return nil
		}
		out, err := os.Create(df)
		if err != nil {
			in.Close()
			log.Warn("Skipping file copy %s: %s", fname, err)
			return nil
		}
		out.Chmod(fi.Mode())
		if _, err = io.Copy(out, in); err != nil {
			log.Warn("Copy from %s to %s failed: %s", fname, df, err)
		}

		if err := out.Close(); err != nil {
			log.Warn("Failed to close %q: %s", df, err)
		}
		if err := in.Close(); err != nil {
			log.Warn("Failed to close reader %q: %s", fname, err)
		}

		return nil
	}
	filepath.Walk(src, walker)
	return failure
}

//CopyFile copies file from src to dst
func CopyFile(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err == nil {
		err = ioutil.WriteFile(dst, data, 0644)
	}
	return err
}
