package manifest

import (
	"testing"
)

var testchart = "../../charts/alpine"

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
