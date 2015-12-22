package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helm/helm/log"
	"github.com/helm/helm/test"
	"github.com/helm/helm/util"
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

func TestByName(t *testing.T) {
	cfg, err := Load("../testdata/Configfile.yaml")
	if err != nil {
		t.Fatalf("Could not load ../testdata/Configfile.yaml: %s", err)
	}

	if rn := cfg.Repos.ByName("charts"); rn != "https://github.com/helm/charts" {
		t.Errorf("Unexpected chart URL: %s", rn)
	}
}

func TestByRepo(t *testing.T) {
	cfg, err := Load("../testdata/Configfile.yaml")
	if err != nil {
		t.Fatalf("Could not load ../testdata/Configfile.yaml: %s", err)
	}

	if rn := cfg.Repos.ByRepo("https://github.com/helm/charts"); rn != "charts" {
		t.Errorf("Unexpected chart name: %s", rn)
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

func TestCanonicalRepo(t *testing.T) {
	expect := "example.com/foo/bar"
	orig := []string{
		"git@example.com:foo/bar.git",
		"git@example.com:/foo/bar.git",
		"http://example.com/foo/bar.git",
		"https://example.com/foo/bar.git",
		"ssh://git@example.com/foo/bar.git",
		"http://example.com/foo/bar",
		"example.com/foo/bar",
	}
	for _, v := range orig {
		cv, err := CanonicalRepo(v)
		if err != nil {
			t.Errorf("Failed to parse %s: %s", v, err)
		}
		if cv != expect {
			t.Errorf("Expected %q, got %q for %q", expect, cv, v)
		}
	}

	expect = "localhost/slurm/bar"
	orig = []string{
		"file:///slurm/bar.git",
		"/slurm/bar.git",
		"localhost/slurm/bar.git",
	}
	for _, v := range orig {
		cv, err := CanonicalRepo(v)
		if err != nil {
			t.Errorf("Failed to parse %s: %s", v, err)
		}
		if cv != expect {
			t.Errorf("Expected %q, got %q for %q", expect, cv, v)
		}
	}
}
