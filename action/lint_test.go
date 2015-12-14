package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helm/helm/test"
	"github.com/helm/helm/util"

	"gopkg.in/yaml.v2"
)

func TestLintSuccess(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "goodChart"

	Create(chartName, tmpHome)

	output := test.CaptureOutput(func() {
		Lint(chartName, tmpHome)
	})

	expected := "Chart [goodChart] has passed all necessary checks"

	test.ExpectContains(t, output, expected)
}

func TestLintMissingReadme(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	os.Remove(filepath.Join(tmpHome, util.WorkspaceChartPath, chartName, "README.md"))

	output := test.CaptureOutput(func() {
		Lint(chartName, tmpHome)
	})

	expected := "A README file was not found"

	if !strings.Contains(output, expected) {
		t.Fatalf("Expected: '%s' in %s.", expected, output)
	}
}

func TestLintMissingChartYaml(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	os.Remove(filepath.Join(tmpHome, util.WorkspaceChartPath, chartName, "Chart.yaml"))

	output := test.CaptureOutput(func() {
		Lint(chartName, tmpHome)
	})

	test.ExpectContains(t, output, "A Chart.yaml file was not found")
	test.ExpectContains(t, output, "Chart [badChart] failed some checks")
}

func TestLintEmptyChartYaml(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	badChartYaml, _ := yaml.Marshal(make(map[string]string))

	chartYaml := filepath.Join(tmpHome, util.WorkspaceChartPath, chartName, "Chart.yaml")

	os.Remove(chartYaml)
	ioutil.WriteFile(chartYaml, badChartYaml, 0644)

	output := test.CaptureOutput(func() {
		Lint(chartName, tmpHome)
	})

	test.ExpectContains(t, output, "Missing Name specification in Chart.yaml file")
	test.ExpectContains(t, output, "Missing Version specification in Chart.yaml file")
	test.ExpectContains(t, output, "Missing description in Chart.yaml file")
	test.ExpectContains(t, output, "Missing maintainers information in Chart.yaml file")
	test.ExpectContains(t, output, fmt.Sprintf("Chart [%s] failed some checks", chartName))
}
