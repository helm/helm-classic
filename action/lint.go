package action

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/helm/helm/log"
	"github.com/helm/helm/util"
	"github.com/helm/helm/validation"
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
	cv := new(validation.ChartValidation)
	cv.Path = chartPath

	//TODO: chartPresenceValidation := v.AddError("Chart found", func(path string, v *ChartValidation) bool { }

	chartYamlValidation := cv.AddError("Chart.yaml is present", func(path string, v *validation.Validation) bool {
		stat, err := os.Stat(v.ChartYamlPath())

		return err == nil && stat.Mode().IsRegular()
	})

	chartYamlNameValidation := chartYamlValidation.AddError("Chart.yaml has a name field", func(path string, v *validation.Validation) bool {
		chartfile, err := v.Chartfile()

		return err == nil && chartfile.Name != ""
	})

	chartYamlNameValidation.AddError("Name declared in Chart.yaml is the same as chart name.", func(path string, v *validation.Validation) bool {
		chartfile, err := v.Chartfile()

		return err == nil && chartfile.Name == cv.ChartName()

	})

	chartYamlValidation.AddError("Chart.yaml has a version field", func(path string, v *validation.Validation) bool {
		chartfile, err := v.Chartfile()

		return err == nil && chartfile.Version != ""
	})

	chartYamlValidation.AddWarning("Chart.yaml has a description field", func(path string, v *validation.Validation) bool {
		chartfile, err := v.Chartfile()

		return err == nil && chartfile.Description != ""
	})

	chartYamlValidation.AddWarning("Chart.yaml has a maintainers field", func(path string, v *validation.Validation) bool {
		chartfile, err := v.Chartfile()

		return err == nil && chartfile.Maintainers != nil
	})

	if cv.Valid() {
		log.Info("Chart[%s] has passed all necessary checks", cv.ChartName())
	} else {
		log.Err("Chart [%s] is not completely valid", cv.ChartName())
	}
}
