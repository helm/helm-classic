package action

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/helm/helm/log"
	"github.com/helm/helm/test"
)

func TestTemplate(t *testing.T) {
	dir := filepath.Join(test.HelmRoot, "testdata/template")
	tpl := filepath.Join(dir, "one.tpl")
	val := filepath.Join(dir, "one.toml")
	var out bytes.Buffer
	o := log.Stdout
	log.Stdout = &out
	defer func() { log.Stdout = o }()

	// TOML
	Template("", tpl, val)
	if out.String() != "Hello World!\n" {
		t.Errorf("Expected Hello World!, got %q", out.String())
	}

	// YAML
	val = filepath.Join(dir, "one.yaml")
	out.Reset()
	Template("", tpl, val)
	if out.String() != "Hello World!\n" {
		t.Errorf("Expected Hello World!, got %q", out.String())
	}

	// JSON
	val = filepath.Join(dir, "one.json")
	out.Reset()
	Template("", tpl, val)
	if out.String() != "Hello World!\n" {
		t.Errorf("Expected Hello World!, got %q", out.String())
	}

	// No data
	out.Reset()
	Template("", tpl, "")
	if out.String() != "Hello Clowns!\n" {
		t.Errorf("Expected Hello Clowns!, got %q", out.String())
	}
}
