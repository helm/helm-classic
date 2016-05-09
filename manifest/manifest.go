// Package manifest provides tools for working with Kubernetes manifests.
package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/helm/helm-classic/codec"
	"github.com/helm/helm-classic/log"
)

// Files gets a list of all manifest files inside of a chart.
//
// chartDir should contain the path to a chart (the directory which
// holds a Chart.yaml file).
//
// This returns an error if it can't access the directory.
func Files(chartDir string) ([]string, error) {
	dir := filepath.Join(chartDir, "manifests")
	files := []string{}

	if _, err := os.Stat(dir); err != nil {
		return files, err
	}

	// add manifest files
	walker := func(fname string, fi os.FileInfo, e error) error {
		if e != nil {
			log.Warn("Encountered error walking %q: %s", fname, e)
			return nil
		}

		if fi.IsDir() {
			return nil
		}

		if filepath.Ext(fname) == ".yaml" || filepath.Ext(fname) == ".yml" {
			files = append(files, fname)
		}

		return nil
	}
	filepath.Walk(dir, walker)

	return files, nil
}

// Manifest represents a Kubernetes manifest object.
type Manifest struct {
	Version, Kind, Name string
	// Filename of source, "" if no file
	Source          string
	VersionedObject *codec.Object
}

// Parse takes a filename, loads the file, and parses it into one or more *Manifest objects.
func Parse(filename string) ([]*Manifest, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ms := []*Manifest{}

	docs, err := codec.YAML.Decode(d).All()
	if err != nil {
		return ms, fmt.Errorf("%s: %s", filename, err)
	}

	for _, doc := range docs {
		ref, err := doc.Meta()
		if err != nil {
			return nil, fmt.Errorf("%s: %s", filename, err)
		}

		m := &Manifest{Version: ref.APIVersion, Kind: ref.Kind, Name: ref.Name, VersionedObject: doc, Source: filename}
		ms = append(ms, m)
	}
	return ms, nil
}

// ParseDir parses all of the manifests inside of a chart directory.
//
// The directory should be the Chart directory (contains Chart.yaml and manifests/)
//
// This will return an error if the directory does not exist, or if there is an
// error parsing or decoding any yaml files.
func ParseDir(chartDir string) ([]*Manifest, error) {
	dir := filepath.Join(chartDir, "manifests")
	files := []*Manifest{}

	if _, err := os.Stat(dir); err != nil {
		return files, err
	}

	// add manifest files
	walker := func(fname string, fi os.FileInfo, e error) error {
		log.Debug("Parsing %s", fname)
		// Chauncey was right.
		if e != nil {
			return e
		}

		if fi.IsDir() {
			return nil
		}

		if filepath.Ext(fname) != ".yaml" && filepath.Ext(fname) != ".yml" {
			log.Debug("Skipping %s. Not a YAML file.", fname)
			return nil
		}

		m, err := Parse(fname)
		if err != nil {
			return err
		}

		files = append(files, m...)

		return nil
	}

	return files, filepath.Walk(dir, walker)
}

// IsKeeper returns true if a manifest has a "helm-keep": "true" annotation.
func IsKeeper(data []byte) bool {
	// Look for ("helm-keep": "true") up to 10 lines after ("annotations:").
	var keep = regexp.MustCompile(`\"annotations\":\s+\{(\n.*){0,10}\"helm-keep\":\s+\"true\"`)
	return keep.Match(data)
}
