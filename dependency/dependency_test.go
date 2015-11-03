package dependency

import (
	"path/filepath"
	"testing"

	"github.com/deis/helm/chart"
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
