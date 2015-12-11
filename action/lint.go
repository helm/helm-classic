package action

import (
	"fmt"
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

// LintAll vlaidates all charts are well-formed
//
// - homedir is the home directory for the user
func LintAll(homedir string) {
	md := filepath.Join(homedir, WorkspaceChartPath, "*")
	charts, err := filepath.Glob(md)
	if err != nil {
		log.Warn("Could not find any charts in %q: %s", md, err)
	}

	if len(charts) == 0 {
		log.Warn("Could not find any charts in %q", md)
	} else {
		for _, c := range charts {
			Lint(filepath.Base(c), homedir)
		}
	}
}

// Lint validates that a chart is well-formed
//
// - chartName is the name of the chart
// - homedir is the home directory for the user
func Lint(chartName, homedir string) {
	var errors = make([]string, 0)

	//assumes chart is in your workspace directory
	chartPath := filepath.Join(homedir, WorkspaceChartPath, chartName)
	//makes sure all files are in place
	structure, fatalDs := directoryStructure(chartPath)

	if len(fatalDs) == 0 {
		dsErrs := checkDirectoryStructure(structure, chartPath)
		errors = append(errors, dsErrs...)
	} else {
		errors = append(errors, fatalDs...)
	}

	//checks to see if chart name is unique
	nameErr := verifyChartNameUnique(chartName)
	if nameErr == nil {
		errors = append(errors, fmt.Sprintf("Chart name %s already exists in charts repository [github.com/helm/charts]. If you're planning on submitting this chart to the charts repo, please consider changing the chart name.", chartName))
	}
	errors = append(errors, verifyMetadata(chartPath)...)
	errors = append(errors, verifyManifests(chartPath)...)

	if len(errors) > 0 {
		for _, message := range errors {
			log.Err(message)
		}
		log.Warn("Chart [%s] failed some checks", chartName)
	} else {
		log.Info("Chart [%s] has passed all necessary checks", chartName)
	}
}

func directoryStructure(chartPath string) (map[string]os.FileInfo, []string) {
	var messages = make([]string, 0)
	structure := make(map[string]os.FileInfo)

	chartInfo, err := os.Stat(chartPath)
	if err != nil {
		messages = append(messages, fmt.Sprintf("Chart %s not found in workspace. Error: %v", chartPath, err))
	}

	if chartInfo.IsDir() {
		files, _ := ioutil.ReadDir(chartPath)
		for _, f := range files {
			structure[f.Name()] = f
		}
	} else {
		messages = append(messages, fmt.Sprintf("Chart Path [%s] is not a directory.", chartPath))
	}

	return structure, messages
}

func checkDirectoryStructure(structure map[string]os.FileInfo, chartPath string) []string {
	var messages = make([]string, 0)

	if _, ok := structure["README.md"]; ok != true {
		messages = append(messages, fmt.Sprintf("A README file was not found in %s", chartPath))
	}

	if _, ok := structure["Chart.yaml"]; ok != true {
		messages = append(messages, fmt.Sprintf("A Chart.yaml file was not found in %s", chartPath))
	}

	manifestInfo, ok := structure["manifests"]

	if ok && manifestInfo.IsDir() {
		// manifest files logic
	} else {
		messages = append(messages, fmt.Sprintf("A manifests directory was not found in %s", chartPath))
	}

	return messages
}

// verifyMetadata checks the Chart.yaml file for a Name, Version, Description, and Maintainers
func verifyMetadata(chartPath string) []string {
	var errors = make([]string, 0)
	var y *chart.Chartfile

	file := filepath.Join(chartPath, "Chart.yaml")
	b, err := ioutil.ReadFile(file)

	if err != nil {
		return append(errors, fmt.Sprint(err))
	}
	if err = yaml.Unmarshal(b, &y); err != nil {
		return append(errors, fmt.Sprint(err))
	}
	//require name, version, description, maintaners
	if y.Name == "" {
		errors = append(errors, "Missing Name specification in Chart.yaml file")
	}
	if y.Version == "" {
		errors = append(errors, "Missing Version specification in Chart.yaml file")
	}
	if y.Description == "" {
		errors = append(errors, "Missing description in Chart.yaml file")
	}
	if y.Maintainers == nil {
		errors = append(errors, "Missing maintainers information in Chart.yaml file")
	}

	return errors
}

func verifyManifests(chartPath string) []string {
	var errors = make([]string, 0)
	manifests, err := manifest.ParseDir(chartPath)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Error walking manifest files. Err: %s", err))
	}

	for _, m := range manifests {
		meta, _ := m.VersionedObject.Meta()
		if meta.Name == "" {
			errors = append(errors, fmt.Sprintf("missing name in %s", m.Source))
		}

		val, ok := meta.Labels["heritage"]
		if !ok || (val != "helm") {
			errors = append(errors, fmt.Sprintf("Missing a label: `heritage: helm` in %s", m.Source))
		}

		kind := meta.Kind
		validKinds := InstallOrder
		valid := validKind(kind, validKinds)
		if !valid {
			errors = append(errors, fmt.Sprintf("%s is not a valid `kind` value for manifest. Here are valid kinds of manifests: %v", kind, validKinds))
		}
	}

	return errors
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
