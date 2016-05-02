package action

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/helm/helm-classic/kubectl"
	"github.com/helm/helm-classic/test"
)

func TestInstall(t *testing.T) {
	// Todo: add tests
	// - with an invalid chart name
	// - with failure to check dependencies
	// - with failure to check dependencies and force option
	// - with chart in current directly
	tests := []struct {
		name     string // Todo: print name on fail
		chart    string
		force    bool
		expected []string
		client   kubectl.Runner
	}{
		{
			name:     "with valid input",
			chart:    "redis",
			expected: []string{"hello from redis"},
			client: TestRunner{
				out: []byte("hello from redis"),
			},
		},
		{
			name:     "with dry-run option",
			chart:    "redis",
			expected: []string{"[CMD] kubectl create -f -"},
			client:   kubectl.PrintRunner{},
		},
		{
			name:     "with unsatisfied dependencies",
			chart:    "kitchensink",
			expected: []string{"Stopping install. Re-run with --force to install anyway."},
			client:   TestRunner{},
		},
		{
			name:     "with unsatisfied dependencies and force option",
			chart:    "kitchensink",
			force:    true,
			expected: []string{"Unsatisfied dependencies", "Running `kubectl create -f`"},
			client:   TestRunner{},
		},
		{
			name:     "with a kubectl error",
			chart:    "redis",
			expected: []string{"Failed to upload manifests: oh snap"},
			client: TestRunner{
				err: errors.New("oh snap"),
			},
		},
	}

	tmpHome := test.CreateTmpHome()
	defer os.RemoveAll(tmpHome)
	test.FakeUpdate(tmpHome)

	// Todo: get rid of this hacky mess
	pp := os.Getenv("PATH")
	defer os.Setenv("PATH", pp)
	os.Setenv("PATH", filepath.Join(test.HelmRoot, "testdata")+":"+pp)

	for _, tt := range tests {
		actual := test.CaptureOutput(func() {
			Install(tt.chart, tmpHome, "", tt.force, false, []string{}, tt.client)
		})

		for _, exp := range tt.expected {
			test.ExpectContains(t, actual, exp)
		}
	}
}
