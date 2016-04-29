package search

import (
	"strings"
	"testing"

	"github.com/helm/helm-classic/config"
)

var testConfig = &config.Configfile{
	Repos: &config.Repos{
		Default: "charts",
		Tables: []*config.Table{
			{Name: "charts", Repo: "https://github.com/helm/charts"},
		},
	},
}

func TestSortScore(t *testing.T) {
	in := []*Result{
		{Name: "bbb", Score: 0},
		{Name: "aaa", Score: 5},
		{Name: "abb", Score: 5},
		{Name: "aab", Score: 0},
		{Name: "bab", Score: 5},
	}
	expect := []string{"aab", "bbb", "aaa", "abb", "bab"}
	expectScore := []int{0, 0, 5, 5, 5}
	SortScore(in)

	// Test Score
	for i := 0; i < len(expectScore); i++ {
		if expectScore[i] != in[i].Score {
			t.Errorf("Sort error on index %d: expected %d, got %d", i, expectScore[i], in[i].Score)
		}
	}
	// Test Name
	for i := 0; i < len(expect); i++ {
		if expect[i] != in[i].Name {
			t.Errorf("Sort error: expected %s, got %s", expect[i], in[i].Name)
		}
	}
}

var testCacheDir = "../testdata/"

func TestSearchByName(t *testing.T) {
	term := "homeslice"

	i := NewIndex(testConfig, testCacheDir)
	charts, err := i.Search(term, 100, false)
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
	term := "homeskillet"

	i := NewIndex(testConfig, testCacheDir)
	charts, err := i.Search(term, 100, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(charts) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(charts))
	}

	for _, v := range charts {
		ch, _ := i.Chart(v.Name)
		if ch.Description != term {
			t.Fatalf("Expected result description to match %s, got %s", term, ch.Description)
		}
	}
}

func TestSearchRegexp(t *testing.T) {
	term := "home[a-z]+"
	i := NewIndex(testConfig, testCacheDir)
	charts, err := i.Search(term, 100, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(charts) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(charts))
	}
}

func TestSearchNotFound(t *testing.T) {
	term := "nonexistent"

	i := NewIndex(testConfig, testCacheDir)
	charts, err := i.Search(term, 100, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(charts) > 0 {
		t.Fatalf("Expected 0 results, got %d", len(charts))
	}
}

func TestCalcScore(t *testing.T) {
	i := NewIndex(testConfig, testCacheDir)

	fields := []string{"aaa", "bbb", "ccc", "ddd"}
	matchline := strings.Join(fields, sep)
	if r := i.calcScore(2, matchline); r != 0 {
		t.Errorf("Expected 0, got %d", r)
	}
	if r := i.calcScore(5, matchline); r != 1 {
		t.Errorf("Expected 1, got %d", r)
	}
	if r := i.calcScore(10, matchline); r != 2 {
		t.Errorf("Expected 2, got %d", r)
	}
	if r := i.calcScore(14, matchline); r != 3 {
		t.Errorf("Expected 3, got %d", r)
	}
}
