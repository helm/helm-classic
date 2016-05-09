package action

import (
	"errors"
	"os"
	"testing"

	"github.com/helm/helm-classic/kubectl"
	"github.com/helm/helm-classic/test"
)

func TestUninstall(t *testing.T) {
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
		{
			name:  "with a helmc annotation",
			chart: "keep",
			force: true,
			expected: []string{"Running `kubectl delete` ...", "Not uninstalling",
				"because of \"helm-keep\" annotation"},
		},
	}

	tmpHome := test.CreateTmpHome()
	defer os.RemoveAll(tmpHome)
	test.FakeUpdate(tmpHome)

	for _, tt := range tests {
		Fetch(tt.chart, "", tmpHome)

		actual := test.CaptureOutput(func() {
			Uninstall(tt.chart, tmpHome, "default", tt.force, tt.client)
		})

		for _, exp := range tt.expected {
			test.ExpectContains(t, actual, exp)
		}
	}
}
