package action

import (
	"os"
	"testing"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/test"
	"github.com/helm/helm/util"
)

func TestFetch(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	defer os.RemoveAll(tmpHome)
	test.FakeUpdate(tmpHome)

	chartFilePath := util.CacheDirectory(tmpHome, "charts/redis/Chart.yaml")
	cf, err := chart.LoadChartfile(chartFilePath)
	if err != nil {
		t.Fatalf("Failed to load chartfile %s: %s", chartFilePath, err)
	}

	Fetch("redis", "", tmpHome, false)

	// Assert that the expected directories exist.
	for _, dep := range cf.Dependencies {
		loc := util.WorkspaceChartDirectory(tmpHome, dep.Name)
		if _, err := os.Stat(loc); err != nil {
			t.Errorf("Expected chart %s: %s", loc, err)
		}
	}
}
