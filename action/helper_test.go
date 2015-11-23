package action

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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

func containsStr(t *testing.T, a, b string) {
	if !strings.Contains(a, b) {
		t.Errorf("\n[Expected]\n%s\n[To Contain]\n%s\n", a, b)
	}
}
