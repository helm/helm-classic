package parameters

import (
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	p, err := Load("../testdata/parameters/someParams.yaml")
	if err != nil {
		t.Fatalf("Could not load ../testdata/parameters/someParams.yaml: %s", err)
	}
	if p.Values["wine"] != "Shiraz" {
		t.Fatalf("Expected wine: Shiraz but got %s", p.Values["wine"])
	}

	if err := p.Save("../testdata/parameters/someParams.yaml-SAVE.yaml"); err != nil {
		t.Fatalf("Could not save: %s", err)
	}

	if _, err := os.Stat("../testdata/parameters/someParams.yaml-SAVE.yaml"); err != nil {
		t.Fatalf("Saved file does not exist: %s", err)
	}

	if err := os.Remove("../testdata/parameters/someParams.yaml-SAVE.yaml"); err != nil {
		t.Fatalf("Could not remove file: %s", err)
	}
}

func TestSaveChartParameters(t *testing.T) {
	chartName := "mychart"
	valueFolder := "../testdata/tmp-chart-values"
	p := &Parameters{Values: map[string]string{
		"wine": "Merlot",
		"beer": "Morretti",
	}}
	err := SaveChartParameters(valueFolder, chartName, p)
	if err != nil {
		t.Fatalf("Could not save chart parameters: %s", err)
	}
	p2, err := LoadChartParameters(valueFolder, chartName)
	if err != nil {
		t.Fatalf("Could not load chart parameters: %s", err)
	}
	if p2.Values["wine"] != p.Values["wine"] {
		t.Fatalf("Loaded chart parameters had the wrong wine %s value when expecting %s", p2.Values["wine"], p.Values["wine"])
	}
	filename := chartParametersFilename(valueFolder, chartName)
	if err := os.Remove(filename); err != nil {
		t.Fatalf("Could not remove file: %s", err)
	}
}

