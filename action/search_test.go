package action

import (
	"testing"

	"github.com/helm/helm-classic/test"
)

func TestSearch(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	Search("homeslice", tmpHome, false)
}

func TestSearchNotFound(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	// test that a "no chart found" message was printed
	expected := "No results found"

	actual := test.CaptureOutput(func() {
		Search("nonexistent", tmpHome, false)
	})

	test.ExpectContains(t, actual, expected)
}
