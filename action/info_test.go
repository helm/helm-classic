package action

import (
	"bytes"
	"os"
	"testing"

	"github.com/helm/helm/log"
	"github.com/helm/helm/test"
)

func TestInfo(t *testing.T) {

	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	var output bytes.Buffer
	log.Stdout = &output

	format := ""

	Info("kitchensink", tmpHome, format)

	expected := `Name: kitchensink
Home: http://github.com/helm/helm
Version: 0.0.1
Description: All the things, all semantically, none working
Details: This package provides a sampling of all of the different manifest types. It can be used to test ordering and other properties of a chart.
`
	test.ExpectEquals(t, output.String(), expected)

	// reset log
	log.Stdout = os.Stdout
}

func TestInfoFormat(t *testing.T) {

	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	var output bytes.Buffer
	log.Stdout = &output

	format := `Hello {{.Name}}`

	Info("kitchensink", tmpHome, format)

	expected := `Hello kitchensink`

	test.ExpectEquals(t, output.String(), expected)

	// reset log
	log.Stdout = os.Stdout
}
