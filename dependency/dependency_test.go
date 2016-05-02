package dependency

import (
	"path/filepath"
	"testing"

	"github.com/helm/helm-classic/chart"
)

var testInstalldir = "../testdata/charts"

func init() {
	var err error
	testInstalldir, err = filepath.Abs(testInstalldir)
	if err != nil {
		panic(err)
	}
}

func TestResolve(t *testing.T) {
	cf, err := chart.LoadChartfile(filepath.Join(testInstalldir, "deptest/Chart.yaml"))
	if err != nil {
		t.Errorf("Could not load chartfile deptest/Chart.yaml: %s", err)
	}

	missed, err := Resolve(cf, testInstalldir)
	if err != nil {
		t.Errorf("could not resolve deps in %s: %s", testInstalldir, err)
	}
	if len(missed) != 2 {
		t.Fatalf("Expected dep3 and honkIfYouLoveDucks to be returned")
	}

	if missed[0].Name != "dep3" {
		t.Errorf("Expected dep3 in slot 0. Got %s", missed[0].Name)
	}
	if missed[1].Name != "honkIfYouLoveDucks" {
		t.Errorf("Expected honkIfYouLoveDucks in slot 1. Got %s", missed[1].Name)
	}
}

func TestOptRepoMatch(t *testing.T) {
	a := &chart.Dependency{
		Repo: "git@github.com:drink/slurm.git",
	}
	b := &chart.Dependency{
		Repo: "https://github.com/drink/slurm.git",
	}

	if !optRepoMatch(a, b) {
		t.Errorf("Expected %s to match %s", a.Repo, b.Repo)
	}

	if !optRepoMatch(a, &chart.Dependency{}) {
		t.Errorf("Expected empty required repo to match any filled from repo.")
	}

	if optRepoMatch(&chart.Dependency{}, a) {
		t.Errorf("Expected empty from repo to fail for a non-empty required repo.")
	}
}

func TestCanonicalRepo(t *testing.T) {
	expect := "example.com/foo/bar.git"
	orig := []string{
		"git@example.com:foo/bar.git",
		"git@example.com:/foo/bar.git",
		"http://example.com/foo/bar.git",
		"https://example.com/foo/bar.git",
		"ssh://git@example.com/foo/bar.git",
	}
	for _, v := range orig {
		cv, err := canonicalRepo(v)
		if err != nil {
			t.Errorf("Failed to parse %s: %s", v, err)
		}
		if cv != expect {
			t.Errorf("Expected %q, got %q for %q", expect, cv, v)
		}
	}

	expect = "localhost/slurm/bar.git"
	orig = []string{
		"file:///slurm/bar.git",
		"/slurm/bar.git",
		"slurm/bar.git",
	}
	for _, v := range orig {
		cv, err := canonicalRepo(v)
		if err != nil {
			t.Errorf("Failed to parse %s: %s", v, err)
		}
		if cv != expect {
			t.Errorf("Expected %q, got %q for %q", expect, cv, v)
		}
	}
}

func TestSatisfies(t *testing.T) {
	a := &chart.Dependency{
		Name:    "slurm",
		Version: "1.2.3",
		Repo:    "ssh://git@example.com/drink/slurm.git",
	}
	b := &chart.Dependency{
		Name:    "slurm",
		Version: "~1.2",
		Repo:    "https://example.com/drink/slurm.git",
	}
	aa := &chart.Dependency{
		Name:    "slurm",
		Version: "1.3.5",
		Repo:    "https://example.com/drink/slurm.git",
	}
	if !satisfies(a, b) {
		t.Errorf("Expected a to satisfy b")
	}

	if satisfies(aa, b) {
		t.Errorf("Expected aa to not satisfy b because of version constraint.")
	}
}
