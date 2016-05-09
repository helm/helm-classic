package manifest

import (
	"io/ioutil"
	"strings"
	"testing"
)

var testchart = "../testdata/charts/kitchensink"
var testPlainManifest = "../testdata/service.json"
var testKeeperManifest = "../testdata/service-keep.json"

func TestFiles(t *testing.T) {
	fs, err := Files(testchart)
	if err != nil {
		t.Errorf("Failed to open %s: %s", testchart, err)
	}

	if len(fs) == 0 {
		t.Errorf("Expected at least one manifest file")
	}
}

func TestParse(t *testing.T) {

	files, _ := Files(testchart)
	found := 0
	for _, file := range files {
		if _, err := Parse(file); err != nil {
			t.Errorf("Failed to parse %s: %s", file, err)
		}
		found++
	}

	if found == 0 {
		t.Errorf("Found no manifests to test.")
	}

}

func TestParseDir(t *testing.T) {
	manifests, err := ParseDir(testchart)
	if err != nil {
		t.Errorf("Failed to parse dir %s: %s", testchart, err)
	}

	target, _ := Files(testchart)
	if len(manifests) < len(target) {
		t.Errorf("Expected at least %d manifests. Got %d", len(target), len(manifests))
	}

	for _, man := range manifests {
		if man.Source == "" {
			t.Error("No file set in manifest.Source.")
		}
		if man.Kind == "" {
			t.Error("Expected kind")
		}
		if man.Name == "" {
			t.Error("Expected name")
		}
	}

	// now test parsing bad files in a chart!
	testchart = "../testdata/charts/malformed"
	manifests, err = ParseDir(testchart)
	if err == nil {
		t.Errorf("Failed to raise an error when parsing dir %s", testchart)
	}
	if !strings.Contains(err.Error(), "malformed.yaml") {
		t.Errorf("Failed to identify which manifest failed to be parsed. Got %s", err)
	}
}

func TestIsKeeper(t *testing.T) {
	// test that an ordinary JSON manifest doesn't look like a keeper
	data, err := ioutil.ReadFile(testPlainManifest)
	if err != nil {
		t.Error(err)
	}
	if IsKeeper(data) {
		t.Errorf("Expected false for %s", testPlainManifest)
	}

	// test that a keeper JSON manifest is detected
	data, err = ioutil.ReadFile(testKeeperManifest)
	if err != nil {
		t.Error(err)
	}
	if !IsKeeper(data) {
		t.Errorf("Expected true for %s", testKeeperManifest)
	}
}
