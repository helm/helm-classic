package action

import (
	"bytes"
	"os"
	"testing"

	"github.com/deis/helm/log"
)

func TestInfo(t *testing.T) {

	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

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
	expect(t, output.String(), expected)

	// reset log
	log.Stdout = os.Stdout
}

func TestInfoFormat(t *testing.T) {

	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	var output bytes.Buffer
	log.Stdout = &output

	format := `Hello {{.Name}}`

	Info("kitchensink", tmpHome, format)

	expected := `Hello kitchensink`

	expect(t, output.String(), expected)

	// reset log
	log.Stdout = os.Stdout
}
