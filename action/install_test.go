package action

import (
	"bytes"
	_ "errors"
	"os"
	"strings"
	"testing"

	"github.com/helm/helm/kubectl"
	"github.com/helm/helm/log"
)

func captureInstallOut(chartName, home, ns string, force, dryRun bool) string {
	var output bytes.Buffer

	log.Stdout = &output
	log.Stderr = &output

	Install(chartName, home, ns, force, dryRun)

	return strings.TrimSpace(output.String())
}

func TestInstall(t *testing.T) {
	tests := []struct {
		chartName string
		force     bool
		dryRun    bool
		expected  []string
		runner    kubectl.Runner
	}{
		{
			"redis",
			false,
			false,
			[]string{"hello from redis"},
			TestRunner{
				out: []byte("hello from redis"),
			},
		},
		// with dry-run set
		{
			"redis",
			false,
			true,
			[]string{"Performing a dry run of `kubectl create -f`"},
			TestRunner{},
		},
		//  with unsatisfied dependencies
		//{
		//"kitchensink",
		//false,
		//false,
		//"Performing a dry run of `kubectl create -f`",
		//TestRunner{},
		//},
		//  with unsatisfied dependencies and force set
		{
			"kitchensink",
			true,
			false,
			[]string{"Unsatisfied dependencies", "Running `kubectl create -f`"},
			TestRunner{},
		},
		// with kubectl error
		//{
		//"redis",
		//false,
		//false,
		//"Failed to upload manifests",
		//TestRunner{
		//err: errors.New("oh snap"),
		//},
		//},
	}

	tmpHome := createTmpHome()
	defer os.RemoveAll(tmpHome)
	fakeUpdate(tmpHome)

	for _, tt := range tests {
		Kubectl = tt.runner
		actual := captureInstallOut(tt.chartName, tmpHome, "", tt.force, tt.dryRun)

		for _, exp := range tt.expected {
			containsStr(t, actual, exp)
		}
	}
}
