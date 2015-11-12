// Package manifest provides tools for working with Kubernetes manifests.
package manifest

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/deis/helm/log"
	"k8s.io/kubernetes/pkg/api"
	_ "k8s.io/kubernetes/pkg/api/v1" // side-effect imports required to enable k8s APIs
	_ "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/runtime"
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
	Source          string
	VersionedObject interface{}
}

// Parse takes a filename, loads the file, and parses it into one or more *Manifest objects.
func Parse(filename string) ([]*Manifest, error) {
	in, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	ms := []*Manifest{}

	docs, err := SplitYAML(in)
	in.Close()
	if err != nil {
		return ms, err
	}

	for _, doc := range docs {
		data, err := yaml.ToJSON(doc)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse %s: %s", filename, err)
		}

		vo, version, kind, err := api.Scheme.Raw().DecodeToVersionedObject(data)
		if err != nil {
			return ms, err
		}

		m := &Manifest{Version: version, Kind: kind, VersionedObject: vo, Source: filename}
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

// SplitYAML splits a compount YAML file into an array of byte arrays.
//
// Each byte array contains an entire YAML "file".
func SplitYAML(r io.Reader) ([][]byte, error) {
	docs := [][]byte{}

	scanner := bufio.NewScanner(r)
	scanner.Split(SplitYAMLDocument)

	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) > 0 {
			docs = append(docs, b)
		}
	}
	return docs, scanner.Err()
}

const yamlSeparator = "\n---"

// SplitYAMLDocument is a bufio.SplitFunc for splitting a YAML document into individual documents.
//
// This is from Kubernetes' 'pkg/util/yaml'.splitYAMLDocument, which is unfortunately
// not exported.
func SplitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// MarshalJSON encodes data with Kubernetes' versioned API.
func MarshalJSON(obj interface{}, version string) ([]byte, error) {
	o, ok := obj.(runtime.Object)
	if !ok {
		return nil, errors.New("Not an Object")
	}
	return api.Scheme.EncodeToVersion(o, version)
}
