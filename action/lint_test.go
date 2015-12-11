package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLintSuccess(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	chartName := "goodChart"

	Create(chartName, tmpHome)

	output := capture(func() {
		Lint(chartName, tmpHome)
	})

	expected := "Chart [goodChart] has passed all necessary checks"

	if !strings.Contains(output, expected) {
		t.Fatalf("Expected: '%s' in %s.", expected, output)
	}
}

func TestLintMissingReadme(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	os.Remove(filepath.Join(tmpHome, WorkspaceChartPath, chartName, "README.md"))

	output := capture(func() {
		Lint(chartName, tmpHome)
	})

	expected := "A README file was not found"

	if !strings.Contains(output, expected) {
		t.Fatalf("Expected: '%s' in %s.", expected, output)
	}
}

func TestLintMissingChartYaml(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	os.Remove(filepath.Join(tmpHome, WorkspaceChartPath, chartName, "Chart.yaml"))

	output := capture(func() {
		Lint(chartName, tmpHome)
	})

	expectContains(t, output, "A Chart.yaml file was not found")
	expectContains(t, output, "Chart [badChart] failed some checks")
}

func TestLintAllNone(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	output := capture(func() {
		Cli().Run([]string{"helm", "--home", tmpHome, "lint", "--all"})
	})

	expectContains(t, output, fmt.Sprintf("Could not find any charts in \"%s", tmpHome))
}

func TestLintEmptyChartYaml(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	badChartYaml, _ := yaml.Marshal(make(map[string]string))

	chartYaml := filepath.Join(tmpHome, WorkspaceChartPath, chartName, "Chart.yaml")

	os.Remove(chartYaml)
	ioutil.WriteFile(chartYaml, badChartYaml, 0644)

	output := capture(func() {
		Lint(chartName, tmpHome)
	})

	expectContains(t, output, "Missing Name specification in Chart.yaml file")
	expectContains(t, output, "Missing Version specification in Chart.yaml file")
	expectContains(t, output, "Missing description in Chart.yaml file")
	expectContains(t, output, "Missing maintainers information in Chart.yaml file")
	expectContains(t, output, fmt.Sprintf("Chart [%s] failed some checks", chartName))
}

func TestLintAll(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	missingReadmeChart := "missingReadme"

	Create(missingReadmeChart, tmpHome)
	os.Remove(filepath.Join(tmpHome, WorkspaceChartPath, missingReadmeChart, "README.md"))

	Create("goodChart", tmpHome)

	output := capture(func() {
		Cli().Run([]string{"helm", "--home", tmpHome, "lint", "--all"})
	})

	expectMatches(t, output, "A README file was not found.*"+missingReadmeChart)
	expectContains(t, output, "Chart [goodChart] has passed all necessary checks")
	expectContains(t, output, "Chart [missingReadme] failed some checks")
}
