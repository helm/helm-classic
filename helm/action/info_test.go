package action

import (
	"bytes"
	"testing"

	"github.com/deis/helm/helm/log"
)

func TestInfo(t *testing.T) {

	tmpHome := createTmpHome()

	var output bytes.Buffer
	log.Stdout = &output

	format := ""

	Info("kitchensink", tmpHome, format)

	expected := `Name: kitchensink
Home: http://github.com/deis/helm
Version: 0.0.1
Description: All the things, all semantically, none working
Details: This package provides a sampling of all of the different manifest types. It can be used to test ordering and other properties of a chart.
`

	if output.String() != expected {
		t.Errorf("Expected %v - Got %v ", expected, output.String())
	}
}

func TestInfoFormat(t *testing.T) {

	tmpHome := createTmpHome()

	var output bytes.Buffer
	log.Stdout = &output

	format := `Hello {{.Name}}`

	Info("kitchensink", tmpHome, format)

	expected := `Hello kitchensink`

	if output.String() != expected {
		t.Errorf("Expected %v - Got %v ", expected, output.String())
	}
}
