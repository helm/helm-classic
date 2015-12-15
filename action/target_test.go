package action

import (
	"testing"

	"github.com/helm/helm/test"
)

func TestTarget(t *testing.T) {
	client := TestRunner{
		out: []byte("lookin good"),
	}

	expected := "lookin good"

	actual := test.CaptureOutput(func() {
		Target(client)
	})

	test.ExpectContains(t, actual, expected)
}
