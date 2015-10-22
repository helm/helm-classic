// Package manifest provides tools for working with Kubernetes manifests.
package manifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/deis/helm/helm/log"
	"k8s.io/kubernetes/pkg/api"
	_ "k8s.io/kubernetes/pkg/api/v1"
	_ "k8s.io/kubernetes/pkg/apis/experimental"
	"k8s.io/kubernetes/pkg/util/yaml"
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

		if filepath.Ext(fname) == ".yaml" {
			files = append(files, fname)
		}

		return nil
	}
	filepath.Walk(dir, walker)

	return files, nil
}

// Manifest represents a Kubernetes manifest object.
type Manifest struct {
	Version, Kind   string
	VersionedObject interface{}
}

// Parse takes a filename, loads the file, and parses it into a *Manifest.
func Parse(filename string) (*Manifest, error) {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data, err := yaml.ToJSON(in)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %s", filename, err)
	}

	vo, version, kind, err := api.Scheme.Raw().DecodeToVersionedObject(data)
	if err != nil {
		return nil, err
	}

	return &Manifest{Version: version, Kind: kind, VersionedObject: vo}, err
}
