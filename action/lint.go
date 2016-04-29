package action

import (
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/github"
	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/manifest"
	"github.com/helm/helm-classic/util"
	"github.com/helm/helm-classic/validation"
)

const (
	// Owner is default Helm repository owner or organization.
	Owner = "helm"

	// Project is the default Charts repository name.
	Project = "charts"

	// MaxMetadataNameLength is the longest Metadata.name allowed by kubernetes.
	MaxMetadataNameLength = 24
)

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

	chartPresenceValidation := cv.AddError("Chart found at "+chartPath, func(path string, v *validation.Validation) bool {
		stat, err := os.Stat(chartPath)
		cv.Path = chartPath

		return err == nil && stat.Mode().IsDir()
	})

	chartYamlPresenceValidation := chartPresenceValidation.AddError("Chart.yaml is present", func(path string, v *validation.Validation) bool {
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

	chartYamlNameValidation.AddError("Name declared in Chart.yaml is the same as directory name.", func(path string, v *validation.Validation) bool {
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

	chartPresenceValidation.AddWarning("README.md is present and not empty", func(path string, v *validation.Validation) bool {
		readmePath := filepath.Join(path, "README.md")
		stat, err := os.Stat(readmePath)

		return err == nil && stat.Mode().IsRegular() && stat.Size() > 0
	})

	manifestsValidation := chartPresenceValidation.AddError("Manifests directory is present", func(path string, v *validation.Validation) bool {
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

	manifestsParsingValidation.AddWarning("Manifests have correct and valid metadata", func(path string, v *validation.Validation) bool {

		success := true
		validKinds := InstallOrder

		for _, m := range cv.Manifests {
			meta, _ := m.VersionedObject.Meta()
			if meta.Name == "" || len(meta.Name) > MaxMetadataNameLength {
				success = false
			}

			if match, _ := regexp.MatchString(`[a-z]([-a-z0-9]*[a-z0-9])?`, meta.Name); !match {
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
		log.Info("Chart [%s] has passed all necessary checks", cv.ChartName())
	} else {
		if cv.ErrorCount > 0 {
			log.Err("Chart [%s] has failed some necessary checks. Check out the error and warning messages listed.", cv.ChartName())
		} else {
			log.Warn("Chart [%s] has passed all necessary checks but failed some checks as well. Proceed with caution. Check out the warnings listed.", cv.ChartName())
		}
	}
}
