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

func TestUninstall(t *testing.T) {
	var output bytes.Buffer
	log.Stdout = &output
	log.Stderr = &output
	defer func() {
		log.Stdout = os.Stdout
		log.Stderr = os.Stderr
	}()

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
			force:    true,
			expected: []string{"Running `kubectl delete` ...", "hello from redis"},
			client: TestRunner{
				out: []byte("hello from redis"),
			},
		},
		{
			name:     "with a kubectl error",
			chart:    "redis",
			force:    true,
			expected: []string{"Running `kubectl delete` ...", "Could not delete Pod redis (Skipping): oh snap"},
			client: TestRunner{
				err: errors.New("oh snap"),
			},
		},
	}

	tmpHome := createTmpHome()
	defer os.RemoveAll(tmpHome)
	fakeUpdate(tmpHome)

	for _, tt := range tests {
		Fetch(tt.chart, "", tmpHome)

		Uninstall(tt.chart, tmpHome, "", tt.force, tt.client)
		actual := output.String()

		for _, exp := range tt.expected {
			if !strings.Contains(actual, exp) {
				t.Errorf("\n[Expected]\n%s\n[To Contain]\n%s\n", actual, exp)
			}
		}
		output.Reset()
	}
}
