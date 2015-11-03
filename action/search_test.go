package action

import (
	"testing"
)

func TestSearchByName(t *testing.T) {
	tmpHome := createTmpHome()
	term := "homeslice"

	charts, err := searchAll(term, tmpHome)
	if err != nil {
		t.Fatal(err)
	}

	if len(charts) != 1 {
		t.Fatalf("Expected 2 results, got %d", len(charts))
	}

	for _, chart := range charts {
		if chart.Name != term {
			t.Fatalf("Expected result name to match %s, got %s", term, chart.Name)
		}
	}
}

func TestSearchByDescription(t *testing.T) {
	tmpHome := createTmpHome()
	term := "homeskillet"

	charts, err := searchAll(term, tmpHome)
	if err != nil {
		t.Fatal(err)
	}

	if len(charts) != 1 {
		t.Fatalf("Expected 2 results, got %d", len(charts))
	}

	for _, chart := range charts {
		if chart.Description != term {
			t.Fatalf("Expected result description to match %s, got %s", term, chart.Description)
		}
	}
}

func TestSearch(t *testing.T) {
	tmpHome := createTmpHome()
	term := "homeslice"

	Search(term, tmpHome)
}
