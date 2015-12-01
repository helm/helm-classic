package action

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/helm/helm/log"
)

func TestSearch(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	term := "homeslice"

	Search(term, tmpHome, false)
}

func TestSearchNotFound(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	term := "nonexistent"

	// capture stdout for testing
	old := log.Stderr
	r, w, _ := os.Pipe()
	log.Stderr = w

	Search(term, tmpHome, false)

	// read test output and restore previous stdout
	w.Close()
	out, _ := ioutil.ReadAll(r)
	log.Stderr = old
	output := string(out[:])

	// test that a "no chart found" message was printed
	txt := "No results found"
	if !strings.Contains(output, txt) {
		t.Fatalf("Expected %s to be printed, got %s", txt, output)
	}
}
