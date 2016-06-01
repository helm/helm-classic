package action

import (
	"os"
	"testing"

	"github.com/helm/helm-classic/test"
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
		{"kitchensink", mockNotFoundGetter, false, "All clear! You have successfully removed kitchensink from your workspace."},

		// when manifests are installed
		{"kitchensink", mockFoundGetter, false, "Found 12 installed manifests for kitchensink.  To remove a chart that has been installed the --force flag must be set."},

		// when manifests are installed and force is set
		{"kitchensink", mockNotFoundGetter, true, "All clear! You have successfully removed kitchensink from your workspace."},

		// when kubectl cannot connect
		{"kitchensink", mockFailConnection, false, "Could not determine if kitchensink is installed.  To remove the chart --force flag must be set."},

		// when kubectl cannot connect and force is set
		{"kitchensink", mockFailConnection, true, "All clear! You have successfully removed kitchensink from your workspace."},
	}

	for _, tt := range tests {
		tmpHome := test.CreateTmpHome()
		test.FakeUpdate(tmpHome)

		Fetch("kitchensink", "", tmpHome)

		// set the mock getter
		kubeGet = tt.getter

		actual := test.CaptureOutput(func() {
			Remove(tt.chartName, tmpHome, tt.force)
		})

		test.ExpectContains(t, actual, tt.expected)

		os.Remove(tmpHome)
	}
}
