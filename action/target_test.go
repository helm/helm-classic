package action

import (
	"bytes"
	"strings"
	"testing"

	"github.com/helm/helm/log"
)

func captureTargetOut() string {
	var output bytes.Buffer

	log.Stdout = &output
	log.Stderr = &output

	Target()

	return strings.TrimSpace(output.String())
}

func TestTarget(t *testing.T) {
	Kubectl = TestRunner{
		out: []byte("lookin good"),
	}

	expected := "lookin good"
	actual := captureTargetOut()

	expect(t, actual, expected)
}
