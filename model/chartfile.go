package model

import (
	"io/ioutil"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"
)

// Chart describes a Helm Chart (e.g. Chart.yaml)
type Chartfile struct {
	Name         string            `yaml:"name"`
	Home         string            `yaml:"home"`
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
}

// Load loads a Chart.yaml file into a *Chart.
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
