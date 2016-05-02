package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/helm/helm-classic/test"
	"github.com/helm/helm-classic/util"

	"gopkg.in/yaml.v2"
)

func TestLintSuccess(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "goodChart"

	Create(chartName, tmpHome)

	output := test.CaptureOutput(func() {
		Lint(util.WorkspaceChartDirectory(tmpHome, chartName))
	})

	expected := "Chart [goodChart] has passed all necessary checks"

	test.ExpectContains(t, output, expected)
}

func TestLintMissingReadme(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	os.Remove(filepath.Join(util.WorkspaceChartDirectory(tmpHome, chartName), "README.md"))

	output := test.CaptureOutput(func() {
		Lint(util.WorkspaceChartDirectory(tmpHome, chartName))
	})

	test.ExpectContains(t, output, "README.md is present and not empty : false")
}

func TestLintMissingChartYaml(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	os.Remove(filepath.Join(util.WorkspaceChartDirectory(tmpHome, chartName), Chartfile))

	output := test.CaptureOutput(func() {
		Lint(util.WorkspaceChartDirectory(tmpHome, chartName))
	})

	test.ExpectContains(t, output, "Chart.yaml is present : false")
	test.ExpectContains(t, output, "Chart [badChart] has failed some necessary checks.")
}

func TestLintMismatchedChartNameAndDir(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	chartName := "chart-0"
	chartDir := "chart-1"
	chart := newSkelChartfile(chartName)
	createWithChart(chart, chartDir, tmpHome)

	output := test.CaptureOutput(func() {
		Lint(util.WorkspaceChartDirectory(tmpHome, chartDir))
	})

	test.ExpectContains(t, output, "Name declared in Chart.yaml is the same as directory name. : false")
}

func TestLintMissingManifestDirectory(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "brokeChart"

	Create(chartName, tmpHome)

	os.RemoveAll(filepath.Join(util.WorkspaceChartDirectory(tmpHome, chartName), "manifests"))

	output := test.CaptureOutput(func() {
		Lint(util.WorkspaceChartDirectory(tmpHome, chartName))
	})

	test.ExpectMatches(t, output, "Manifests directory is present : false")
	test.ExpectContains(t, output, "Chart ["+chartName+"] has failed some necessary checks")
}

func TestLintEmptyChartYaml(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "badChart"

	Create(chartName, tmpHome)

	badChartYaml, _ := yaml.Marshal(make(map[string]string))

	chartYaml := util.WorkspaceChartDirectory(tmpHome, chartName, Chartfile)

	os.Remove(chartYaml)
	ioutil.WriteFile(chartYaml, badChartYaml, 0644)

	output := test.CaptureOutput(func() {
		Lint(util.WorkspaceChartDirectory(tmpHome, chartName))
	})

	test.ExpectContains(t, output, "Chart.yaml has a name field : false")
	test.ExpectContains(t, output, "Chart.yaml has a version field : false")
	test.ExpectContains(t, output, "Chart.yaml has a description field : false")
	test.ExpectContains(t, output, "Chart.yaml has a maintainers field : false")
	test.ExpectContains(t, output, fmt.Sprintf("Chart [%s] has failed some necessary checks", chartName))
}

func TestLintBadPath(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	chartName := "badChart"

	output := test.CaptureOutput(func() {
		Lint(util.WorkspaceChartDirectory(tmpHome, chartName))
	})

	msg := "Chart found at " + tmpHome + "/workspace/charts/" + chartName + " : false"
	test.ExpectContains(t, output, msg)
}
