package action

import (
	"strings"
	"testing"

	"github.com/helm/helm/test"
)

func TestSearch(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	Search("homeslice", tmpHome, false)
}

func TestSearchNotFound(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	output := test.CaptureOutput(func() {
		Search("nonexistent", tmpHome, false)
	})

	// test that a "no chart found" message was printed
	txt := "No results found"

	if !strings.Contains(output, txt) {
		t.Fatalf("Expected %s to be printed, got %s", txt, output)
	}
}
