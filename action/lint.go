package action

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/helm/helm/log"
	"github.com/helm/helm/manifest"
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

	chartYamlPresenceValidation := cv.AddError("Chart.yaml is present", func(path string, v *validation.Validation) bool {
		stat, err := os.Stat(v.ChartYamlPath())

		return err == nil && stat.Mode().IsRegular()
	})

	chartYamlValidation := chartYamlPresenceValidation.AddError("Chart.yaml is valid yaml", func(path string, v *validation.Validation) bool {
		chartfile, err := v.Chartfile()
		if err == nil {
			cv.Chartfile = chartfile
		}

		return err == nil
	})

	chartYamlNameValidation := chartYamlValidation.AddError("Chart.yaml has a name field", func(path string, v *validation.Validation) bool {
		return cv.Chartfile.Name != ""
	})

	chartYamlNameValidation.AddError("Name declared in Chart.yaml is the same as chart name.", func(path string, v *validation.Validation) bool {
		return cv.Chartfile.Name == cv.ChartName()
	})

	chartYamlValidation.AddError("Chart.yaml has a version field", func(path string, v *validation.Validation) bool {
		return cv.Chartfile.Version != ""
	})

	chartYamlValidation.AddWarning("Chart.yaml has a description field", func(path string, v *validation.Validation) bool {
		return cv.Chartfile.Description != ""
	})

	chartYamlValidation.AddWarning("Chart.yaml has a maintainers field", func(path string, v *validation.Validation) bool {
		return cv.Chartfile.Maintainers != nil
	})

	cv.AddWarning("README.md is present", func(path string, v *validation.Validation) bool {
		readmePath := filepath.Join(path, "README.md")
		stat, err := os.Stat(readmePath)

		return err == nil && stat.Mode().IsRegular()
	})

	manifestsValidation := cv.AddError("Manifests directory is present", func(path string, v *validation.Validation) bool {
		stat, err := os.Stat(v.ChartManifestsPath())

		return err == nil && stat.Mode().IsDir()
	})

	manifestsParsingValidation := manifestsValidation.AddError("Manifests are valid yaml", func(path string, v *validation.Validation) bool {
		manifests, err := manifest.ParseDir(cv.Path)
		if err == nil {
			cv.Manifests = manifests
		}

		return err == nil && cv.Manifests != nil
	})

	manifestsParsingValidation.AddError("Manifests have correct and valid metadata", func(path string, v *validation.Validation) bool {

		success := true
		validKinds := InstallOrder

		for _, m := range cv.Manifests {
			meta, _ := m.VersionedObject.Meta()
			if meta.Name == "" {
				success = false
			}

			val, ok := meta.Labels["heritage"]
			if !ok || (val != "helm") {
				success = false
			}

			kind := meta.Kind
			validManifestKind := false

			for _, validKind := range validKinds {
				if kind == validKind {
					validManifestKind = true
				}
			}

			if validManifestKind == false {
				success = false
			}
		}

		return success
	})

	if cv.Valid() {
		log.Info("Chart[%s] has passed all necessary checks", cv.ChartName())
	} else {
		log.Err("Chart [%s] is not completely valid", cv.ChartName())
	}
}
