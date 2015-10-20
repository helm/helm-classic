package model

import (
	"testing"
)

const testfile = "../testdata/test-Chart.yaml"

func TestLoad(t *testing.T) {
	f, err := Load(testfile)
	if err != nil {
		t.Errorf("Error loading %s: %s", testfile, err)
	}

	if f.Name != "alpine-pod" {
		t.Errorf("Expected alpine-pod, got %s", f.Name)
	}

	if len(f.Maintainers) != 2 {
		t.Errorf("Expected 2 maintainers, got %d", len(f.Maintainers))
	}

	if len(f.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(f.Dependencies))
	}

	if f.Dependencies[1].Name != "bar" {
		t.Errorf("Expected second dependency to be bar: %q", f.Dependencies[1].Name)
	}
}

func TestVersionOK(t *testing.T) {
	f, err := Load(testfile)
	if err != nil {
		t.Errorf("Error loading %s: %s", testfile, err)
	}

	// These are canaries. The SemVer package exhuastively tests the
	// various  permutations. This will alert us if we wired it up
	// incorrectly.

	d := f.Dependencies[1]
	if d.VersionOK("1.0.0") {
		t.Errorf("1.0.0 should have been marked out of range")
	}

	if !d.VersionOK("1.2.3") {
		t.Errorf("Version 1.2.3 should have been marked in-range")
	}

}
