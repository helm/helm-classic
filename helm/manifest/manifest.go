// Package manifest provides tools for working with Kubernetes manifests.
package manifest

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/deis/helm/helm/log"
	//"k8s.io/kubernetes/pkg/api"
	//_ "k8s.io/kubernetes/pkg/api/v1"
	//_ "k8s.io/kubernetes/pkg/apis/experimental"
	//"k8s.io/kubernetes/pkg/runtime"
	//"k8s.io/kubernetes/pkg/util/yaml"
	"github.com/technosophos/kubelite/codec"
	//"github.com/technosophos/kubelite/v1"
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
	Version, Kind string
	// Filename of source, "" if no file
	Source     string
	Definition *codec.Manifest
}

// Parse takes a filename, loads the file, and parses it into one or more *Manifest objects.
func Parse(filename string) ([]*Manifest, error) {
	ms := []*Manifest{}

	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return ms, err
	}

	// Parse all of the manifests in the project
	ys, err := codec.YAML.Decode(d).All()
	if err != nil {
		return ms, err
	}
	ms = make([]*Manifest, len(ys))
	for i, y := range ys {
		r, err := y.Ref()
		if err != nil {
			log.Warn("Skip. Error fetching reference: %s", err)
			continue
		}

		ms[i] = &Manifest{
			Kind:       r.Kind,
			Version:    r.APIVersion,
			Source:     filename,
			Definition: y,
		}
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

		if filepath.Ext(fname) != ".yaml" {
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

func MarshalJSON(v interface{}, version string) ([]byte, error) {
	var b bytes.Buffer
	err := codec.JSON.Encode(&b).One(v)
	log.Info("Generated JSON: %s", b.String())
	return b.Bytes(), err
}
