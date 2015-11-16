package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTmpHome() string {
	tmpHome, _ := ioutil.TempDir("", "helm_home")
	defer os.Remove(tmpHome)
	return tmpHome
}

func TestEnsureRepo(t *testing.T) {
	tmpHome := createTmpHome()

	repo := "https://github.com/helm/charts"
	ensureRepo(repo, filepath.Join(tmpHome, "cache", "charts"))
}

func TestParseConfigfile(t *testing.T) {
	cfg, err := Parse([]byte(DefaultConfigfile))
	if err != nil {
		t.Fatalf("Could not parse DefaultConfigfile: %s", err)
	}

	r := cfg.Repos
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

func TestLoadConfigfile(t *testing.T) {
	cfg, err := Load("../testdata/Configfile.yaml")
	if err != nil {
		t.Fatalf("Could not load ../testdata/Configfile.yaml: %s", err)
	}

	if len(cfg.Repos.Tables) != 3 {
		t.Errorf("Expected 3 remotes.")
	}
}

func TestSave(t *testing.T) {
	cfg, err := Load("../testdata/Configfile.yaml")
	if err != nil {
		t.Fatalf("Could not load ../testdata/Configfile.yaml: %s", err)
	}

	if err := cfg.Save("../testdata/Configfile-SAVE.yaml"); err != nil {
		t.Fatalf("Could not save: %s", err)
	}

	if _, err := os.Stat("../testdata/Configfile-SAVE.yaml"); err != nil {
		t.Fatalf("Saved file does not exist: %s", err)
	}

	if err := os.Remove("../testdata/Configfile-SAVE.yaml"); err != nil {
		t.Fatalf("Could not remove file: %s", err)
	}
}
