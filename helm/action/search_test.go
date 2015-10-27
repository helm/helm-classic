package action

import (
	"testing"
)

func TestSanitizeTerm(t *testing.T) {
	if s := sanitizeTerm(""); s != "*" {
		t.Errorf("Expected *, got %s", s)
	}

	if s := sanitizeTerm("foo"); s != "*foo*" {
		t.Errorf("Expected *foo*, got %q", s)
	}
}

func TestSearch(t *testing.T) {
	home := "../testdata/"
	term := "earchte"

	lines, err := search(term, home)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(lines))
	}
}
