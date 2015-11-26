package parameters

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/helm/helm/log"

	"gopkg.in/yaml.v2"
)

// Parameters is a per Chart configuration file to keep track of things like environment specific properties.
type Parameters struct {
	Values     map[string]string     `yaml:"values"`
}

// LoadChartParameters loads parameter values for the given chart name and param folder
// if the file exists or returns an empty object if not
func LoadChartParameters(paramFolder string, chartName string) (*Parameters, error) {
	filename := chartParametersFilename(paramFolder, chartName)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &Parameters{Values: make(map[string]string)}, nil
	}
	log.Debug("Loading chart %s parameters from %s", chartName, filename)
	return Load(filename)
}

// SaveChartParameters saves the overriden properties for the given chart name so that they
// can be used next the the chart is installed or updated
func SaveChartParameters(paramFolder string, chartName string, customParams *Parameters) error {
	folderName := filepath.Join(paramFolder, chartName)
	err := os.MkdirAll(folderName, 0755)
	if err != nil {
		return err
	}
	filename := chartParametersFilename(paramFolder, chartName)
	log.Debug("Saving chart %s parameters to %s", chartName, filename)
	return customParams.Save(filename)
}

func chartParametersFilename(paramFolder string, chartName string) string {
	return filepath.Join(paramFolder, chartName, "parameters.yaml")
}

// Load loads parameter values by filename.
func Load(filename string) (*Parameters, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	p, err := Parse(b)
	if err != nil {
		return p, err
	}
	return p, nil
}

// Parse parses a byte slice into a *Parameters.
func Parse(data []byte) (*Parameters, error) {
	p := &Parameters{}
	if err := yaml.Unmarshal(data, p); err != nil {
		return p, err
	}
	return p, nil
}

// Save writes the Parameters as YAML into the named file.
func (p *Parameters) Save(filename string) error {
	b, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0755)
}

