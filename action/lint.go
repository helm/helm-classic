package action

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/helm/helm/chart"
	"github.com/helm/helm/log"
	"github.com/helm/helm/manifest"
	"gopkg.in/yaml.v2"
)

// Owner is default Helm repository owner or organization.
var Owner = "helm"

// Project is the default Charts repository name.
var Project = "charts"

// RepoService is a GitHub client instance.
var RepoService GHRepoService

// GHRepoService is a restricted interface to GitHub client operations.
type GHRepoService interface {
	DownloadContents(string, string, string, *github.RepositoryContentGetOptions) (io.ReadCloser, error)
}

// Lint validates that a chart is well-formed
//
// - chartName is the name of the chart
// - homedir is the home directory for the user
func Lint(chartName, homedir string) {

	//assumes chart is in your workspace directory
	chartPath := filepath.Join(homedir, WorkspaceChartPath, chartName)

	//makes sure all files are in place
	structure := directoryStructure(chartPath)
	checkDirectoryStructure(structure, chartPath)

	//checks to see if chart name is unique
	nameErr := verifyChartNameUnique(chartName)
	if nameErr == nil {
		log.Warn("Chart name %s already exists in charts repository [github.com/helm/charts]. If you're planning on submitting this chart to the charts repo, please consider changing the chart name.", chartName)
	}

	verifyMetadata(chartPath)
	verifyManifests(chartPath)

	log.Info("Chart [%s] has passed all necessary checks", chartName)
}

func directoryStructure(chartPath string) map[string]os.FileInfo {
	structure := make(map[string]os.FileInfo)

	chartInfo, err := os.Stat(chartPath)
	if err != nil {
		log.Die("Chart %s not found in workspace. Error: %v", chartPath, err)
	}

	if chartInfo.IsDir() {
		files, _ := ioutil.ReadDir(chartPath)
		for _, f := range files {
			structure[f.Name()] = f
		}
	} else {
		log.Die("Chart Path [%s] is not a directory.", chartPath)
	}

	return structure
}

func checkDirectoryStructure(structure map[string]os.FileInfo, chartPath string) {
	if _, ok := structure["README.md"]; ok != true {
		log.Warn("A README file was not found in %s", chartPath)
	}

	if _, ok := structure["Chart.yaml"]; ok != true {
		log.Die("A Chart.yaml file was not found in %s", chartPath)
	}

	manifestInfo, ok := structure["manifests"]
	if ok && manifestInfo.IsDir() {
		// manifest files logic
	} else {
		log.Die("A manifests directory was not found in %s", chartPath)
	}
}

// verifyMetadata checks the Chart.yaml file for a Name, Version, Description, and Maintainers
func verifyMetadata(chartPath string) {
	file := filepath.Join(chartPath, "Chart.yaml")
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Die("Error reading Chart.yaml.\nError: ", err)
	}
	var y *chart.Chartfile
	if err = yaml.Unmarshal(b, &y); err != nil {
		log.Die("Error parsing Chart.yaml file. \nError: ", err)
	}
	//require name, version, description, maintaners
	if y.Name == "" {
		log.Die("Missing Name specification in Chart.yaml file")
	}
	if y.Version == "" {
		log.Die("Missing Version specification in Chart.yaml file")
	}
	if y.Description == "" {
		log.Die("Missing description in Chart.yaml file")
	}
	if y.Maintainers == nil {
		log.Die("Missing maintainers information in Chart.yaml file")
	}
}

func verifyManifests(chartPath string) {
	manifests, err := manifest.ParseDir(chartPath)
	if err != nil {
		log.Die("Error walking manifest files. Err: ", err)
	}

	for _, m := range manifests {
		meta, _ := m.VersionedObject.Meta()
		if meta.Name == "" {
			log.Die("Missing Name in %s", m.Source)
		}

		val, ok := meta.Labels["heritage"]
		if !ok || (val != "helm") {
			log.Die("Missing a label: `heritage: helm` in %s", m.Source)
		}

		kind := meta.Kind
		validKinds := InstallOrder
		valid := validKind(kind, validKinds)
		if !valid {
			log.Warn("%s is not a valid `kind` value for manifest. Here are valid kinds of manifests: %v", kind, validKinds)
		}
	}
}

func validKind(kind string, validKinds []string) bool {
	for _, validKind := range validKinds {
		if kind == validKind {
			return true
		}
	}
	return false
}

func verifyChartNameUnique(chartName string) error {
	if RepoService == nil {
		RepoService = github.NewClient(nil).Repositories
	}

	chartPath := filepath.Join(chartName, "Chart.yaml")
	_, err := RepoService.DownloadContents(Owner, Project, chartPath, nil)
	return err
}
