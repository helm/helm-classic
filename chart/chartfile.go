package chart

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/helm/helm-classic/log"
	"gopkg.in/yaml.v2"
)

// Chartfile describes a Helm Chart (e.g. Chart.yaml)
type Chartfile struct {
	Name         string            `yaml:"name"`
	From         *Dependency       `yaml:"from,omitempty"`
	Home         string            `yaml:"home"`
	Source       []string          `yaml:"source,omitempty"`
	Version      string            `yaml:"version"`
	Description  string            `yaml:"description"`
	Maintainers  []string          `yaml:"maintainers,omitempty"`
	Details      string            `yaml:"details,omitempty"`
	Dependencies []*Dependency     `yaml:"dependencies,omitempty"`
	PreInstall   map[string]string `yaml:"preinstall,omitempty"`
}

// Dependency describes a specific dependency.
type Dependency struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Repo    string `yaml:"repo,omitempty"`
}

// LoadChartfile loads a Chart.yaml file into a *Chart.
func LoadChartfile(filename string) (*Chartfile, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var y Chartfile
	return &y, yaml.Unmarshal(b, &y)
}

// Save saves a Chart.yaml file
func (c *Chartfile) Save(filename string) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, b, 0644)
}

// RepoName gets the name of the Git repo, or an empty string if none is found.
func RepoName(chartpath string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Err("Could not get working directory: %s", err)
		return ""
	}
	defer func() {
		if err := os.Chdir(wd); err != nil {
			log.Die("Unrecoverable error: %s", err)
		}
	}()

	if err := os.Chdir(chartpath); err != nil {
		log.Err("Could not find chartpath %s: %s", chartpath, err)
		return ""
	}

	out, err := exec.Command("git", "config", "--get", "remote.origin.url").CombinedOutput()
	if err != nil {
		log.Err("Git failed to get the origin name: %s %s", err, string(out))
		return ""
	}

	return strings.TrimSpace(string(out))
}

// VersionOK returns true if the given version meets the constraints.
//
// It returns false if the version string or constraint is unparsable or if the
// version does not meet the constraint.
func (d *Dependency) VersionOK(version string) bool {
	c, err := semver.NewConstraint(d.Version)
	if err != nil {
		return false
	}
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}

	return c.Check(v)
}
