package action

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/test"
	"github.com/helm/helm-classic/util"
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
	Template("", tpl, val, false)
	if out.String() != "Hello World!\n" {
		t.Errorf("Expected Hello World!, got %q", out.String())
	}

	// force false
	os.Setenv("HELM_FORCE_FLAG", "false")
	if err = Template(tpl, val, "", false); err == nil {
		t.Errorf("Expected error but got nil")
	}
	tpl1 := filepath.Join(dir, "two.yaml")
	util.CopyFile(tpl, tpl1)
	// force true
	if err = Template(tpl1, val, "", true); err != nil {
		t.Errorf("error force-generating template (%s)", err.Error())
	}
	defer os.Remove(tpl1)

	// YAML
	val = filepath.Join(dir, "one.yaml")
	out.Reset()
	Template("", tpl, val, false)
	if out.String() != "Hello World!\n" {
		t.Errorf("Expected Hello World!, got %q", out.String())
	}

	// JSON
	val = filepath.Join(dir, "one.json")
	out.Reset()
	Template("", tpl, val, false)
	if out.String() != "Hello World!\n" {
		t.Errorf("Expected Hello World!, got %q", out.String())
	}

	// No data
	out.Reset()
	Template("", tpl, "", false)
	if out.String() != "Hello Clowns!\n" {
		t.Errorf("Expected Hello Clowns!, got %q", out.String())
	}
}
