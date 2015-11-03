package action

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var helmRoot string

func init() {
	helmRoot = filepath.Join(os.Getenv("GOPATH"), "src/github.com/deis/helm/")
}

func createTmpHome() string {
	// create a temp directory
	tmpHome, _ := ioutil.TempDir("", "helm_home")
	defer os.Remove(tmpHome)

	// create cache directory
	tmpHomeCache := filepath.Join(tmpHome, "cache")
	os.Mkdir(tmpHomeCache, 0755)

	// absolute path to testdata charts
	testChartsPath := filepath.Join(helmRoot, "testdata/charts")

	// copy testdata charts into cache
	// mock git clone
	copyDir(testChartsPath, filepath.Join(tmpHomeCache, "charts"))

	return tmpHome
}
