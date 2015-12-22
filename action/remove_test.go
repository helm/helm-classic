package action

import (
	"os"
	"testing"

	"github.com/helm/helm/test"
)

var (
	// mock responses from kubectl
	mockFoundGetter    = func(string) string { return "installed" }
	mockNotFoundGetter = func(string) string { return "not found" }
	mockFailConnection = func(string) string { return "unable to connect to a server" }
)

func TestTRemove(t *testing.T) {
	kg := kubeGet
	defer func() { kubeGet = kg }()

	tests := []struct {
		chartName string
		getter    kubeGetter
		force     bool
		expected  string
	}{
		{"redis", mockNotFoundGetter, false, "All clear! You have successfully removed redis from your workspace."},

		// when manifests are installed
		{"redis", mockFoundGetter, false, "Found 1 installed manifests for redis.  To remove a chart that has been installed the --force flag must be set."},

		// when manifests are installed and force is set
		{"redis", mockNotFoundGetter, true, "All clear! You have successfully removed redis from your workspace."},

		// when kubectl cannot connect
		{"redis", mockFailConnection, false, "Could not determine if redis is installed.  To remove the chart --force flag must be set."},

		// when kubectl cannot connect and force is set
		{"redis", mockFailConnection, true, "All clear! You have successfully removed redis from your workspace."},
	}

	for _, tt := range tests {
		tmpHome := test.CreateTmpHome()
		test.FakeUpdate(tmpHome)

		Fetch("redis", "", tmpHome, false)

		// set the mock getter
		kubeGet = tt.getter

		actual := test.CaptureOutput(func() {
			Remove(tt.chartName, tmpHome, tt.force)
		})

		test.ExpectContains(t, actual, tt.expected)

		os.Remove(tmpHome)
	}
}
