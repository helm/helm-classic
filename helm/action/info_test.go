package action

import (
	"path/filepath"
	"testing"
)

func TestInfo(t *testing.T) {
	chartPath := filepath.Join(HOME, CacheChartPath, "alpine", "Chart.yaml")

	actual, err := describeChart(chartPath)
	if err != nil {
		t.Error(err)
	}

	expected := `Chart: alpine-pod
Description: Simple pod running Alpine Linux.
Details: This package provides a basic Alpine Linux image that can be used for basic debugging and troubleshooting. By default, it starts up, sleeps for a long time, and then eventually stops.
Version: 0.0.1
Website: http://github.com/deis/helm
From: /Users/adamreese/p/go/src/github.com/deis/helm/helm/testdata/helm_home/cache/charts/alpine/Chart.yaml
Dependencies: []
`
	if actual.String() != expected {
		t.Errorf("Expected %v - Got %v ", expected, actual.String())
	}
}
