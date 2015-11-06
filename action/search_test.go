package action

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/deis/helm/log"
)

func TestSearchByName(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	term := "homeslice"

	charts, err := searchAll(term, tmpHome)
	if err != nil {
		t.Fatal(err)
	}

	if len(charts) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(charts))
	}

	for _, chart := range charts {
		if chart.Name != term {
			t.Fatalf("Expected result name to match %s, got %s", term, chart.Name)
		}
	}
}

func TestSearchByDescription(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	term := "homeskillet"

	charts, err := searchAll(term, tmpHome)
	if err != nil {
		t.Fatal(err)
	}

	if len(charts) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(charts))
	}

	for _, chart := range charts {
		if chart.Description != term {
			t.Fatalf("Expected result description to match %s, got %s", term, chart.Description)
		}
	}
}

func TestSearch(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	term := "homeslice"

	Search(term, tmpHome)
}

func TestSearchNotFound(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	term := "nonexistent"

	// capture stdout for testing
	old := log.Stdout
	r, w, _ := os.Pipe()
	log.Stdout = w

	Search(term, tmpHome)

	// read test output and restore previous stdout
	w.Close()
	out, _ := ioutil.ReadAll(r)
	log.Stdout = old
	output := string(out[:])

	// test that a "no chart found" message was printed
	txt := "No chart found for \"" + term + "\"."
	if !strings.Contains(output, txt) {
		t.Fatalf("Expected %s to be printed, got %s", txt, output)
	}
}
