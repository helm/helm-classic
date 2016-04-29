package action

import (
	"testing"

	"github.com/helm/helm-classic/test"
)

func TestInfo(t *testing.T) {

	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	format := ""
	expected := `Name: kitchensink
Home: http://github.com/helm/helm
Version: 0.0.1
Description: All the things, all semantically, none working
Details: This package provides a sampling of all of the different manifest types. It can be used to test ordering and other properties of a chart.`

	actual := test.CaptureOutput(func() {
		Info("kitchensink", tmpHome, format)
	})

	test.ExpectContains(t, actual, expected)
}

func TestInfoFormat(t *testing.T) {

	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	format := `Hello {{.Name}}`
	expected := `Hello kitchensink`

	actual := test.CaptureOutput(func() {
		Info("kitchensink", tmpHome, format)
	})

	test.ExpectContains(t, actual, expected)
}
