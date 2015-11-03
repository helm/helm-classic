package repo

import (
	"os"
	"testing"
)

func TestParseRepofile(t *testing.T) {
	r, err := ParseRepofile([]byte(DefaultRepofile))
	if err != nil {
		t.Fatalf("Could not parse DefaultRepofile: %s", err)
	}

	if r.Default != "charts" {
		t.Errorf("Expected 'charts', got %q", r.Default)
	}

	if len(r.Tables) != 1 {
		t.Errorf("Expected exactly 1 table.")
	}

	if r.Tables[0].Repo != "https://github.com/deis/charts" {
		t.Errorf("Wrong URL")
	}

	if r.Tables[0].Name != "charts" {
		t.Errorf("Wrong table name")
	}
}

func TestLoadRepofile(t *testing.T) {
	r, err := LoadRepofile("../testdata/Repofile.yaml")
	if err != nil {
		t.Fatalf("Could not load ../testdata/Repofile.yaml: %s", err)
	}

	if len(r.Tables) != 3 {
		t.Errorf("Expected 3 remotes.")
	}
}
func TestSaveRepofile(t *testing.T) {
	r, err := LoadRepofile("../testdata/Repofile.yaml")
	if err != nil {
		t.Fatalf("Could not load ../testdata/Repofile.yaml: %s", err)
	}

	if err := r.Save("../testdata/Repofile-SAVE.yaml"); err != nil {
		t.Fatalf("Could not save: %s", err)
	}

	if _, err := os.Stat("../testdata/Repofile-SAVE.yaml"); err != nil {
		t.Fatalf("Saved file does not exist: %s", err)
	}

	if err := os.Remove("../testdata/Repofile-SAVE.yaml"); err != nil {
		t.Fatalf("Could not remove file: %s", err)
	}
}
