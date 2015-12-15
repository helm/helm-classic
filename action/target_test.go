package action

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/helm/helm/log"
	"github.com/helm/helm/test"
)

func TestTarget(t *testing.T) {
	var output bytes.Buffer
	log.Stdout = &output
	log.Stderr = &output
	defer func() {
		log.Stdout = os.Stdout
		log.Stderr = os.Stderr
	}()

	client := TestRunner{
		out: []byte("lookin good"),
	}

	expected := "lookin good"
	Target(client)
	actual := strings.TrimSpace(output.String())

	test.ExpectEquals(t, actual, expected)
}
