package action

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/helm/helm/log"
)

const tmpConfigfile = `repos:
  default: charts
  tables:
    - name: charts
      repo: https://github.com/helm/charts
`

var helmRoot string

func init() {
	// Turn on debug output, convert os.Exit(1) to panic()
	log.IsDebugging = true

	helmRoot = filepath.Join(os.Getenv("GOPATH"), "src/github.com/helm/helm/")
}

func createTmpHome() string {
	tmpHome, _ := ioutil.TempDir("", "helm_home")
	defer os.Remove(tmpHome)
	return tmpHome
}

func fakeUpdate(home string) {
	ensureHome(home)

	ioutil.WriteFile(filepath.Join(home, Configfile), []byte(tmpConfigfile), 0755)

	// absolute path to testdata charts
	testChartsPath := filepath.Join(helmRoot, "testdata/charts")

	// copy testdata charts into cache
	// mock git clone
	copyDir(testChartsPath, filepath.Join(home, "cache/charts"))
}

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("\n[Expected] type: %v\n%v\n[Got] type: %v\n%v\n", reflect.TypeOf(b), b, reflect.TypeOf(a), a)
	}
}

func expectMatches(t *testing.T, actual string, expected string) {
	regexp := regexp.MustCompile(expected)

	if !regexp.Match([]byte(actual)) {
		t.Errorf("\n[Expected] %v\nto contain %v\n", actual, expected)
	}
}

func expectContains(t *testing.T, actual string, expected string) {
	if !strings.Contains(actual, expected) {
		t.Errorf("\n[Expected] %v\nto contain %v\n", actual, expected)
	}
}

func capture(fn func()) string {
	logStderr := log.Stderr
	logStdout := log.Stdout
	osStdout := os.Stdout
	osStderr := os.Stderr

	r, w, _ := os.Pipe()

	log.Stderr = w
	log.Stdout = w
	os.Stdout = w
	os.Stderr = w

	fn()

	// read test output and restore previous stdout
	w.Close()
	out, _ := ioutil.ReadAll(r)

	defer func() {
		log.Stderr = logStderr
		log.Stdout = logStdout
		os.Stdout = osStdout
		os.Stderr = osStderr
	}()

	return string(out[:])
}
