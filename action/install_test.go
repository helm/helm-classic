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

func captureInstallOut(chartName, home, ns string, force bool, client kubectl.Runner) string {
	var output bytes.Buffer

	log.Stdout = &output
	log.Stderr = &output

	Install(chartName, home, ns, force, client)

	return strings.TrimSpace(output.String())
}

func TestInstall(t *testing.T) {
	tests := []struct {
		chartName string
		force     bool
		expected  []string
		client    kubectl.Runner
	}{
		{
			"redis",
			false,
			[]string{"hello from redis"},
			TestRunner{
				out: []byte("hello from redis"),
			},
		},
		// with dry-run set
		//{
		//"redis",
		//false,
		//[]string{"Performing a dry run of `kubectl create -f`"},
		//TestRunner{},
		//},
		//  with unsatisfied dependencies
		//{
		//"kitchensink",
		//false,
		//"Performing a dry run of `kubectl create -f`",
		//TestRunner{},
		//},
		//  with unsatisfied dependencies and force set
		{
			"kitchensink",
			true,
			[]string{"Unsatisfied dependencies", "Running `kubectl create -f`"},
			TestRunner{},
		},
		// with kubectl error
		//{
		//"redis",
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
		actual := captureInstallOut(tt.chartName, tmpHome, "", tt.force, tt.client)

		for _, exp := range tt.expected {
			containsStr(t, actual, exp)
		}
	}
}
