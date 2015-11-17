package action

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/helm/helm/log"
)

var (
	// save globals for reset
	defaultKubeGet func(string) string
	defaultLogErr  io.Writer
	defaultLogOut  io.Writer

	// output stores command output
	output bytes.Buffer

	// mock responses from kubectl
	mockFoundGetter    = func(string) string { return "installed" }
	mockNotFoundGetter = func(string) string { return "not found" }
	mockFailConnection = func(string) string { return "unable to connect to a server" }
)

func saveDefaults() {
	defaultKubeGet = kubeGet
	defaultLogErr = log.Stderr
	defaultLogOut = log.Stdout
}

func rmCmdTeardown() {
	kubeGet = defaultKubeGet
	log.Stderr = defaultLogErr
	log.Stdout = defaultLogOut
	output.Reset()
}

func TestTRemove(t *testing.T) {
	saveDefaults()

	tests := []struct {
		chartName string
		getter    kubeGetter
		force     bool
		match     string
	}{
		{"kitchensink", mockNotFoundGetter, false, "All clear! You have successfully removed kitchensink from your workspace."},

		// when manifests are installed
		{"kitchensink", mockFoundGetter, false, "Found 8 installed manifests for kitchensink.  To remove a chart that has been installed the --force flag must be set."},

		// when manifests are installed and force is set
		{"kitchensink", mockNotFoundGetter, true, "All clear! You have successfully removed kitchensink from your workspace."},

		// when kubectl cannot connect
		{"kitchensink", mockFailConnection, false, "Could not determine if kitchensink is installed.  To remove the chart --force flag must be set."},

		// when kubectl cannot connect and force is set
		{"kitchensink", mockFailConnection, true, "All clear! You have successfully removed kitchensink from your workspace."},
	}

	for _, tt := range tests {
		tmpHome := createTmpHome()
		fakeUpdate(tmpHome)

		Fetch("kitchensink", "", tmpHome)

		log.Stderr = &output
		log.Stdout = &output

		// set the mock getter
		kubeGet = tt.getter

		Remove(tt.chartName, tmpHome, tt.force)

		if !strings.Contains(output.String(), tt.match) {
			t.Errorf("\nExpected\n%s\nTo contain\n%s\n", output.String(), tt.match)
		}

		os.Remove(tmpHome)
		rmCmdTeardown()
	}
}
