package action

import (
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	Search("homeslice", tmpHome, false)
}

func TestSearchNotFound(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	output := capture(func() {
		Search("nonexistent", tmpHome, false)
	})

	// test that a "no chart found" message was printed
	txt := "No results found"

	if !strings.Contains(output, txt) {
		t.Fatalf("Expected %s to be printed, got %s", txt, output)
	}
}
