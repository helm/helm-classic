package model

import (
	"io/ioutil"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"
)

// Chart describes a Helm Chart (e.g. Chart.yaml)
type Chart struct {
	Name         string        `yaml:"name"`
	Home         string        `yaml:"home"`
	Version      string        `yaml:"version"`
	Description  string        `yaml:"description"`
	Maintainers  []string      `yaml:"maintainers,omitempty"`
	Details      string        `yaml:"details,omitempty"`
	Dependencies []*Dependency `yaml:"dependencies,omitempty"`
	PreInstall   []string      `yaml:"preinstall,omitempty"`
}

// Dependency describes a specific dependency.
type Dependency struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// Load loads a Chart.yaml file into a *Chart.
func Load(filename string) (*Chart, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var y Chart
	return &y, yaml.Unmarshal(b, &y)
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
