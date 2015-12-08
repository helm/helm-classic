package action

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/helm/helm/kubectl"
	"github.com/helm/helm/log"
)

func TestInstall(t *testing.T) {
	var output bytes.Buffer
	log.Stdout = &output
	log.Stderr = &output
	defer func() {
		log.Stdout = os.Stdout
		log.Stderr = os.Stderr
	}()

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

	tmpHome := createTmpHome()
	defer os.RemoveAll(tmpHome)
	fakeUpdate(tmpHome)

	defer func() {
		if err := recover(); err != nil {
			output.WriteString(err.(string))
		}
	}()

	for _, tt := range tests {
		Install(tt.chart, tmpHome, "", tt.force, tt.client)
		actual := output.String()

		for _, exp := range tt.expected {
			if !strings.Contains(actual, exp) {
				t.Logf("%+v", tt)
				t.Errorf("\n[Expected]\n%s\n[To Contain]\n%s\n", actual, exp)
			}
		}
	}
}
