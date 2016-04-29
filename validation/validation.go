package validation

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/helm/helm-classic/chart"
	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/manifest"
	"gopkg.in/yaml.v2"
)

// ChartValidation represents a specific instance of validation against a specific directory.
type ChartValidation struct {
	Path         string
	Validations  []*Validation
	Chartfile    *chart.Chartfile
	Manifests    []*manifest.Manifest
	ErrorCount   int
	WarningCount int
}

const (
	warningLevel = 1
	errorLevel   = 2
)

// Validation represents a single validation of a ChartValidation.
type Validation struct {
	children  []*Validation
	path      string
	validator Validator
	Message   string
	level     int
}

// ChartYamlPath returns the fully qualified path to the "Chart.yaml" file.
func (v *Validation) ChartYamlPath() string {
	return filepath.Join(v.path, "Chart.yaml")
}

// ChartManifestsPath returns the fully qualified path to the "Manifests" directory.
func (v *Validation) ChartManifestsPath() string {
	return filepath.Join(v.path, "manifests")
}

// Chartfile returns a chart.Chartfile object formed from the Chart.yaml file.
func (v *Validation) Chartfile() (*chart.Chartfile, error) {
	var y *chart.Chartfile
	b, err := ioutil.ReadFile(v.ChartYamlPath())

	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(b, &y); err != nil {
		return nil, err
	}

	return y, nil
}

// Validator is a declared function that returns the result of a Validation.
type Validator func(path string, v *Validation) (result bool)

func (cv *ChartValidation) addValidator(v *Validation) {
	cv.Validations = append(cv.Validations, v)
}

func (v *Validation) addValidator(child *Validation) {
	v.children = append(v.children, child)
}

// AddError adds error level validation to a ChartValidation.
func (cv *ChartValidation) AddError(message string, fn Validator) *Validation {
	v := new(Validation)
	v.Message = message
	v.validator = fn
	v.level = errorLevel
	v.path = cv.Path

	cv.addValidator(v)

	return v
}

// AddWarning adds a warning level validation to a ChartValidation
func (cv *ChartValidation) AddWarning(message string, fn Validator) *Validation {
	v := new(Validation)
	v.Message = message
	v.validator = fn
	v.level = warningLevel
	v.path = cv.Path

	cv.addValidator(v)

	return v
}

// AddError adds an error level validation to a Validation.
func (v *Validation) AddError(message string, fn Validator) *Validation {
	child := new(Validation)
	child.Message = message
	child.validator = fn
	child.level = errorLevel
	child.path = v.path

	v.addValidator(child)

	return child
}

// AddWarning adds a warning level validation to a Validation.
func (v *Validation) AddWarning(message string, fn Validator) *Validation {
	child := new(Validation)
	child.Message = message
	child.validator = fn
	child.level = warningLevel
	child.path = v.path

	v.addValidator(child)

	return child
}

// ChartName returns the name of the chart directory.
func (cv *ChartValidation) ChartName() string {
	return filepath.Base(cv.Path)
}

func (v *Validation) valid() bool {
	return v.validator(v.path, v)
}

func (v *Validation) walk(talker func(*Validation) bool) {
	validResult := talker(v)

	if validResult {
		for _, child := range v.children {
			child.walk(talker)
		}
	}
}

func (cv *ChartValidation) walk(talker func(v *Validation) bool) {
	for _, v := range cv.Validations {
		v.walk(talker)
	}
}

// Valid returns true if every validation passes.
func (cv *ChartValidation) Valid() bool {
	var valid = true

	fmt.Printf("\nVerifying %s chart is a valid chart...\n", cv.ChartName())
	cv.walk(func(v *Validation) bool {
		v.path = cv.Path
		vv := v.valid()
		if !vv {
			switch v.level {
			case 2:
				cv.ErrorCount = cv.ErrorCount + 1
				msg := v.Message + " : " + strconv.FormatBool(vv)
				log.Err(msg)
			case 1:
				cv.WarningCount = cv.WarningCount + 1
				msg := v.Message + " : " + strconv.FormatBool(vv)
				log.Warn(msg)
			}
		} else {
			msg := v.Message + " : " + strconv.FormatBool(vv)
			log.Info(msg)
		}

		valid = valid && vv
		return valid
	})

	return valid
}
