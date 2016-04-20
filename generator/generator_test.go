package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkip(t *testing.T) {
	pass := []string{
		".foo/.bar/baz",
		"_bar/foo",
		"/foo/bar/.baz/slurm",
	}
	for _, p := range pass {
		if e := skip(p); e != nil {
			t.Errorf("%s should have been nil, got %v", p, e)
		}
	}

	nopass := []string{
		"foo/.bar",
		"foo/_bar",
		"foo/..",
		"..",
		".",
	}
	for _, f := range nopass {
		if e := skip(f); e != filepath.SkipDir {
			t.Errorf("%s should have been nil, got %v", f, e)
		}
	}
}

func TestReadGenerator(t *testing.T) {
	dir := "../testdata/generator"
	pass := []string{"one.yaml", "two.yaml", "three.txt", "four/four.txt", "four/five.txt"}
	fail := []string{"fail.txt", "fail2.txt"}

	for _, p := range pass {
		f, err := os.Open(filepath.Join(dir, p))
		if err != nil {
			t.Errorf("failed to read %s: %s", p, err)
		}
		out, err := readGenerator(f)
		if err != nil {
			t.Errorf("%s failed read generator: %s", p, err)
		}
		f.Close()

		if out != "echo foo bar baz" {
			t.Errorf("Expected %s to output 'echo foo bar baz', got %q", p, out)
		}
	}
	for _, p := range fail {
		f, err := os.Open(filepath.Join(dir, p))
		if err != nil {
			t.Errorf("failed to read %s: %s", p, err)
		}
		out, err := readGenerator(f)
		if err != nil {
			t.Errorf("%s failed read generator: %s", p, err)
		}
		f.Close()

		if out != "" {
			t.Errorf("Expected %s to output empty string, got %q", p, out)
		}
	}
}

func TestWalk(t *testing.T) {
	dir := "../testdata/generator"
	count, err := Walk(dir, []string{}, false)
	if err != nil {
		t.Fatalf("Failed to walk: %s", err)
	}
	if count != 5 {
		t.Errorf("Expected 5 executes, got %d", count)
	}
}
