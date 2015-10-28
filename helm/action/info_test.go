package action

import (
	"bytes"
	"testing"

	"github.com/deis/helm/helm/log"
)

func TestInfo(t *testing.T) {
	var output bytes.Buffer
	log.Stdout = &output

	format := ""

	Info("alpine", TestHome, format)

	expected := `Name: alpine-pod
Home: http://github.com/deis/helm
Version: 0.0.1
Description: Simple pod running Alpine Linux.
Details: This package provides a basic Alpine Linux image that can be used for basic debugging and troubleshooting. By default, it starts up, sleeps for a long time, and then eventually stops.
`

	if output.String() != expected {
		t.Errorf("Expected %v - Got %v ", expected, output.String())
	}
}

func TestInfoFormat(t *testing.T) {
	var output bytes.Buffer
	log.Stdout = &output

	format := `Hello {{.Name}}`

	Info("alpine", TestHome, format)

	expected := `Hello alpine-pod`

	if output.String() != expected {
		t.Errorf("Expected %v - Got %v ", expected, output.String())
	}
}
