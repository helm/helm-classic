package validation

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/helm/helm/chart"
	"gopkg.in/yaml.v2"
)

// ChartValidation represents a specific instance of validation against a specific directory
type ChartValidation struct {
	Path        string
	Validations []*Validation
}

const (
	warningLevel = 0
	errorLevel   = 1
)

// Validation - hold messages related to validation of something
type Validation struct {
	children  []*Validation
	path      string
	validator validator
	Message   string
	level     int
}

//ChartYamlPath - path to Chart.yaml
func (v *Validation) ChartYamlPath() string {
	return filepath.Join(v.path, "Chart.yaml")
}

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

type validator func(path string, v *Validation) (result bool)

func (cv *ChartValidation) addValidator(v *Validation) {
	cv.Validations = append(cv.Validations, v)
}

func (v *Validation) addValidator(child *Validation) {
	v.children = append(v.children, child)
}

// AddError - add error level validation to ChartValidation
func (cv *ChartValidation) AddError(message string, fn validator) *Validation {
	v := new(Validation)
	v.Message = message
	v.validator = fn
	v.level = errorLevel
	v.path = cv.Path

	cv.addValidator(v)

	return v
}

// AddWarning - add warning level validation to ChartValidation
func (cv *ChartValidation) AddWarning(message string, fn validator) *Validation {
	v := new(Validation)
	v.Message = message
	v.validator = fn
	v.level = warningLevel
	v.path = cv.Path

	cv.addValidator(v)

	return v
}

// AddError - add error level validation to Validation
func (v *Validation) AddError(message string, fn validator) *Validation {
	child := new(Validation)
	child.Message = message
	child.validator = fn
	child.level = errorLevel
	child.path = v.path

	v.addValidator(child)

	return child
}

// AddWarning - add warning level validation to Validation
func (v *Validation) AddWarning(message string, fn validator) *Validation {
	child := new(Validation)
	child.Message = message
	child.validator = fn
	child.level = warningLevel
	child.path = v.path

	v.addValidator(child)

	return child
}

func (cv *ChartValidation) ChartName() string {
	return filepath.Base(cv.Path)
}

func (v *Validation) valid() bool {
	return v.validator(v.path, v)
}

func (v *Validation) walk(talker func(_ *Validation) bool) {
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

// Valid - true if every validation passes
func (cv *ChartValidation) Valid() bool {
	var valid bool = true

	cv.walk(func(v *Validation) bool {
		vv := v.valid()
		fmt.Println(fmt.Sprintf(v.Message+" : %v", vv))
		valid = valid && vv
		return valid
	})

	return valid
}
