package dependency

import (
	"path/filepath"
	"testing"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/config"
	"github.com/helm/helm/util"
)

func TestSatisfies(t *testing.T) {
	cf := &chart.Chartfile{
		Name:    "foo",
		Version: "5.4.6",
	}
	dep1 := &chart.Dependency{
		Name:    "foo",
		Version: ">5.0, <5.5",
	}
	dep2 := &chart.Dependency{
		Name:    "foo",
		Version: "~4",
	}

	if err := Satisfies(cf, dep1); err != nil {
		t.Errorf("Surprise! Got error: %s", err)
	}

	if err := Satisfies(cf, dep2); err == nil {
		t.Errorf("Should have gotten an error because %s does not satisfy %s", cf.Version, dep2.Version)
	}
}

func TestResolver(t *testing.T) {
	home := "../testdata/resolver"
	rpath := filepath.Join(home, util.Configfile)
	cfg, err := config.Load(rpath)
	if err != nil {
		t.Fatal(err.Error())
	}

	c, err := chart.LoadChartfile(util.CacheDirectory(home, "charts/deptest/Chart.yaml"))
	if err != nil {
		t.Fatal(err.Error())
	}

	r := NewResolver(cfg, util.WorkspaceChartDirectory(home), util.CacheDirectory(home))
	deps, err := r.Resolve(c, "charts")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(deps) != 3 {
		t.Errorf("Expected 3 dependencies. Got %d", len(deps))
	}

	for n, d := range deps {
		if !d.Found {
			t.Errorf("Expected %s to be found.", n)
		}
		if d.Fetched {
			t.Errorf("Did not expect %s to be fetched.", n)
		}
		if !d.Satisfies {
			t.Errorf("Expected %s to satisfy requirements. Got: %s", n, d.SatisfiesErr)
		}
	}

}
