package chart

import (
	"testing"
)

const testfile = "../testdata/test-Chart.yaml"
const testchart = "../testdata/charts/kitchensink"
const testtemplatechart = "../testdata/charts/template"

func TestLoad(t *testing.T) {
	c, err := Load(testchart)
	if err != nil {
		t.Errorf("Failed to load chart: %s", err)
	}

	if c.Chartfile.Name != "kitchensink" {
		t.Errorf("Expected chart name to be 'kitchensink'. Got '%s'.", c.Chartfile.Name)
	}
	if c.Chartfile.Dependencies[0].Version != "~10.21" {
		d := c.Chartfile.Dependencies[0].Version
		t.Errorf("Expected dependency 0 to have version '~10.21'. Got '%s'.", d)
	}

	if len(c.Pods) != 3 {
		t.Errorf("Expected 3 pods, got %d", len(c.Pods))
	}

	if _, ok := c.Pods[0].Annotations[OriginFile]; !ok {
		t.Error("Failed to get origin file from pod 0")
	}

	if len(c.ReplicationControllers) == 0 {
		t.Error("No RCs found")
	}
	if _, ok := c.ReplicationControllers[0].Annotations[OriginFile]; !ok {
		t.Error("Failed to get origin file from pod 0")
	}

	if len(c.Namespaces) == 0 {
		t.Errorf("No namespaces found")
	}
	if _, ok := c.Namespaces[0].Annotations[OriginFile]; !ok {
		t.Error("Failed to get origin file from pod 0")
	}

	if len(c.Secrets) == 0 {
		t.Error("Is it secret? Is it safe? NO!")
	}
	if _, ok := c.Secrets[0].Annotations[OriginFile]; !ok {
		t.Error("Failed to get origin file from pod 0")
	}

	if len(c.PersistentVolumes) == 0 {
		t.Errorf("No volumes.")
	}
	if _, ok := c.PersistentVolumes[0].Annotations[OriginFile]; !ok {
		t.Error("Failed to get origin file from pod 0")
	}

	if len(c.Services) == 0 {
		t.Error("No service. Just like [insert mobile provider name here]")
	}
	if _, ok := c.Services[0].Annotations[OriginFile]; !ok {
		t.Error("Failed to get origin file from pod 0")
	}
}

func TestLoadChart(t *testing.T) {
	f, err := LoadChartfile(testfile)
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

	if f.PreInstall["mykeys"] != "generate-keypair foo" {
		t.Errorf("Expected map value for mykeys.")
	}
}

func TestVersionOK(t *testing.T) {
	f, err := LoadChartfile(testfile)
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

func TestLoadTemplate(t *testing.T) {
	c, err := Load(testtemplatechart)
	if err != nil {
		t.Errorf("Failed to load chart: %s", err)
	}

	if c.Chartfile.Name != "template" {
		t.Errorf("Expected chart name to be 'template'. Got '%s'.", c.Chartfile.Name)
	}

	if len(c.Templates) != 1 {
		t.Errorf("Expected 1 templates, got %d", len(c.Templates))
	} else {
		temp := c.Templates[0];
		objects := temp.Objects;
		params := temp.Parameters;

		if len(objects) < 1 {
			t.Errorf("Expected some Objects in the template, got %d", len(objects))
		}
		if len(params) < 1 {
			t.Errorf("Expected some Parameters in the template, got %d", len(params))
		}
	}
}
