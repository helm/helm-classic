package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/test"
	"github.com/helm/helm-classic/util"
)

func TestEnsureRepo(t *testing.T) {
	tmpHome := test.CreateTmpHome()

	repo := "https://github.com/helm/charts"
	ensureRepo(repo, filepath.Join(tmpHome, "cache", "charts"))
}

func TestParseConfigfile(t *testing.T) {
	cfg, err := Parse([]byte(util.DefaultConfigfile))
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

	if r.Tables[0].Repo != "https://github.com/helm/charts" {
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

func TestPrintSummary(t *testing.T) {
	var b bytes.Buffer

	log.Stdout = &b
	log.Stderr = &b
	defer func() {
		log.Stdout = os.Stdout
		log.Stderr = os.Stderr
	}()

	diff := `M	README.md
M	cassandra
A	jenkins
M	mysql
M	owncloud`

	expected := []string{
		"Updated 3 charts\ncassandra                    mysql                        owncloud",
		"Added 1 charts\njenkins",
	}

	printSummary(diff)
	actual := b.String()

	for _, exp := range expected {
		if !strings.Contains(actual, exp) {
			t.Errorf("Expected %q to contain %q", actual, exp)
		}
	}
}
