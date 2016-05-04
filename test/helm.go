package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/util"
)

const tmpConfigfile = `repos:
  default: charts
  tables:
    - name: charts
      repo: https://github.com/helm/charts
`

// HelmRoot - dir root of the project
var HelmRoot = filepath.Join(os.Getenv("GOPATH"), "src/github.com/helm/helm-classic/")

// CreateTmpHome create a temporary directory for $HELMC_HOME
func CreateTmpHome() string {
	tmpHome, _ := ioutil.TempDir("", "helmc_home")
	defer os.Remove(tmpHome)
	return tmpHome
}

// FakeUpdate add testdata to home path
func FakeUpdate(home string) {
	util.EnsureHome(home)

	ioutil.WriteFile(filepath.Join(home, util.Configfile), []byte(tmpConfigfile), 0755)

	// absolute path to testdata charts
	testChartsPath := filepath.Join(HelmRoot, "testdata/charts")

	// copy testdata charts into cache
	// mock git clone
	util.CopyDir(testChartsPath, filepath.Join(home, "cache/charts"))
}

// ExpectEquals assert a == b
func ExpectEquals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("\n[Expected] type: %v\n%v\n[Got] type: %v\n%v\n", reflect.TypeOf(b), b, reflect.TypeOf(a), a)
	}
}

// ExpectMatches assert a ~= b
func ExpectMatches(t *testing.T, actual string, expected string) {
	regexp := regexp.MustCompile(expected)

	if !regexp.Match([]byte(actual)) {
		t.Errorf("\n[Expected] %v\nto contain %v\n", actual, expected)
	}
}

// ExpectContains assert b is contained within a
func ExpectContains(t *testing.T, actual string, expected string) {
	if !strings.Contains(actual, expected) {
		t.Errorf("\n[Expected] %v\nto contain %v\n", actual, expected)
	}
}

// CaptureOutput redirect all log/std streams, capture and replace
func CaptureOutput(fn func()) (out string) {
	logStderr := log.Stderr
	logStdout := log.Stdout
	osStdout := os.Stdout
	osStderr := os.Stderr

	defer func() {
		log.Stderr = logStderr
		log.Stdout = logStdout
		os.Stdout = osStdout
		os.Stderr = osStderr

		if r := recover(); r != nil {
			out = r.(string)
		}
	}()

	r, w, _ := os.Pipe()

	log.Stderr = w
	log.Stdout = w
	os.Stdout = w
	os.Stderr = w

	fn()

	// read test output and restore previous stdout
	w.Close()
	b, _ := ioutil.ReadAll(r)
	out = string(b)
	return
}
