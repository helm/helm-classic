package action

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const tmpRepofile = `default: charts
tables:
  - name: charts
    repo: https://github.com/deis/charts
`

var helmRoot string

func init() {
	helmRoot = filepath.Join(os.Getenv("GOPATH"), "src/github.com/deis/helm/")
}

func createTmpHome() string {
	tmpHome, _ := ioutil.TempDir("", "helm_home")
	defer os.Remove(tmpHome)
	return tmpHome
}

func fakeUpdate(home string) {
	ensureHome(home)

	ioutil.WriteFile(filepath.Join(tmpHomeCache, Repofile), []byte(tmpRepofile), 0755)

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
