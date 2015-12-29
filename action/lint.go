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
	"github.com/helm/helm/util"
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
	md := util.WorkspaceChartDirectory(homedir, "*")
	chartPaths, err := filepath.Glob(md)
	if err != nil {
		log.Warn("Could not find any charts in %q: %s", md, err)
	}

	if len(chartPaths) == 0 {
		log.Warn("Could not find any charts in %q", md)
	} else {
		for _, chartPath := range chartPaths {
			Lint(chartPath)
		}
	}
}

// Lint validates that a chart is well-formed
//
// - chartPath path to chart directory
func Lint(chartPath string) {
	v := new(Validation)
	chartName := filepath.Base(chartPath)

	//makes sure all files are in place
	structure := directoryStructure(chartPath, v)

	if v.Valid() {
		checkDirectoryStructure(structure, chartPath, v)
	}

	//checks to see if chart name is unique
	verifyChartNameUnique(chartName, v)
	verifyMetadata(chartPath, v)
	verifyManifests(chartPath, v)

	numErrors := len(v.Errors)

	if len(v.Warnings) > 0 || numErrors > 0 {
		for _, warning := range v.Warnings {
			log.Warn(warning)
		}
		for _, err := range v.Errors {
			log.Err(err)
		}

		if numErrors > 0 {
			log.Err("Chart [%s] failed some checks", chartName)
		} else {
			log.Warn("Chart [%s] failed some checks", chartName)
		}
	} else {
		log.Info("Chart [%s] has passed all necessary checks", chartName)
	}
}

func directoryStructure(chartPath string, v *Validation) map[string]os.FileInfo {
	structure := make(map[string]os.FileInfo)

	chartInfo, err := os.Stat(chartPath)
	if err != nil {
		v.AddError(fmt.Sprintf("Chart %s not found in workspace. Error: %v", chartPath, err))
		return structure
	}

	if chartInfo.IsDir() {
		files, _ := ioutil.ReadDir(chartPath)
		for _, f := range files {
			structure[f.Name()] = f
		}
	} else {
		v.AddError(fmt.Sprintf("Chart Path [%s] is not a directory.", chartPath))
	}

	return structure
}

func checkDirectoryStructure(structure map[string]os.FileInfo, chartPath string, v *Validation) {
	if _, ok := structure["README.md"]; ok != true {
		v.AddWarning(fmt.Sprintf("A README file was not found in %s", chartPath))
	}

	manifestInfo, ok := structure["manifests"]

	if ok && manifestInfo.IsDir() {
		// manifest files logic
	} else {
		v.AddError(fmt.Sprintf("A manifests directory was not found in %s", chartPath))
	}
}

// verifyMetadata checks the Chart.yaml file for a Name, Version, Description, and Maintainers
func verifyMetadata(chartPath string, v *Validation) {
	var y *chart.Chartfile

	file := filepath.Join(chartPath, Chartfile)
	chartDir := filepath.Base(chartPath)
	b, err := ioutil.ReadFile(file)

	if err != nil {
		v.AddError("A Chart.yaml file was not found")
		return
	}

	if err = yaml.Unmarshal(b, &y); err != nil {
		v.AddError(fmt.Sprint(err))
		return
	}

	// require name, version, description, maintaners
	if y.Name == "" {
		v.AddError("Missing Name specification in Chart.yaml file")
	}
	if y.Name != chartDir {
		v.AddError(fmt.Sprintf("Chart.yaml name (%s) is not the same as its directory (%s)", y.Name, chartDir))
	}
	if y.Version == "" {
		v.AddError("Missing Version specification in Chart.yaml file")
	}
	if y.Description == "" {
		v.AddWarning("Missing description in Chart.yaml file")
	}
	if y.Maintainers == nil {
		v.AddWarning("Missing maintainers information in Chart.yaml file")
	}
}

func verifyManifests(chartPath string, v *Validation) {
	manifests, err := manifest.ParseDir(chartPath)

	if err != nil {
		v.AddError(fmt.Sprintf("Error walking manifest files. Err: %s", err))
	}

	for _, m := range manifests {
		meta, _ := m.VersionedObject.Meta()
		if meta.Name == "" {
			v.AddWarning(fmt.Sprintf("missing name in %s", m.Source))
		}

		val, ok := meta.Labels["heritage"]
		if !ok || (val != "helm") {
			v.AddWarning(fmt.Sprintf("Missing a label: `heritage: helm` in %s", m.Source))
		}

		kind := meta.Kind
		validKinds := InstallOrder

		if valid := validKind(kind, validKinds); !valid {
			v.AddError(fmt.Sprintf("%s is not a valid `kind` value for manifest. Here are valid kinds of manifests: %v", kind, validKinds))
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

func verifyChartNameUnique(chartName string, v *Validation) {
	if RepoService == nil {
		RepoService = github.NewClient(nil).Repositories
	}

	chartPath := filepath.Join(chartName, Chartfile)

	if _, err := RepoService.DownloadContents(Owner, Project, chartPath, nil); err == nil {
		v.AddWarning(fmt.Sprintf("Chart name %s already exists in charts repository [github.com/helm/charts]. If you're planning on submitting this chart to the charts repo, please consider changing the chart name.", chartName))
	}
}
