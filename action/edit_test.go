package action

import (
	"path/filepath"
	"testing"

	"github.com/helm/helm/test"
)

func TestListChart(t *testing.T) {
	chartDir := filepath.Join(test.HelmRoot, "testdata/charts/redis")

	files, err := listChart(chartDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2, got: %v", len(files))
	}

}

func TestJoinChart(t *testing.T) {

	chartDir := filepath.Join(test.HelmRoot, "testdata/charts/redis")

	// prepare files fixture
	var files []string
	paths := []string{
		"testdata/charts/redis/Chart.yaml",
		"testdata/charts/redis/manifests/redis-pod.yaml",
	}
	for _, f := range paths {
		files = append(files, filepath.Join(test.HelmRoot, f))
	}

	bytes, err := joinChart(chartDir, files)
	if err != nil {
		t.Fatalf("failed to join chart: %v", bytes)
	}

	if bytes.Len() < 1 {
		t.Fatalf("empty chart after join: %v", bytes)
	}

}

// TODO: TestSaveChart to test serialization
