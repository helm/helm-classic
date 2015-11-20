package action

import (
	"bytes"
	"strings"
	"testing"

	"github.com/helm/helm/kubectl"
	"github.com/helm/helm/log"
)

func captureTargetOut(client kubectl.Runner) string {
	var output bytes.Buffer

	log.Stdout = &output
	log.Stderr = &output

	Target(client)

	return strings.TrimSpace(output.String())
}

func TestTarget(t *testing.T) {
	client := TestRunner{
		out: []byte("lookin good"),
	}

	expected := "lookin good"
	actual := captureTargetOut(client)

	expect(t, actual, expected)
}
