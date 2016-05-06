package cli

import (
	"fmt"
	"os"
	"testing"

	"github.com/helm/helm-classic/action"
	"github.com/helm/helm-classic/test"
	"github.com/helm/helm-classic/util"
)

func TestLintAllNone(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	output := test.CaptureOutput(func() {
		Cli().Run([]string{"helmc", "--home", tmpHome, "lint", "--all"})
	})

	test.ExpectContains(t, output, fmt.Sprintf("Could not find any charts in \"%s", tmpHome))
}

func TestLintSingle(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "goodChart"
	action.Create(chartName, tmpHome)

	output := test.CaptureOutput(func() {
		Cli().Run([]string{"helmc", "--home", tmpHome, "lint", chartName})
	})

	test.ExpectContains(t, output, fmt.Sprintf("Chart [%s] has passed all necessary checks", chartName))
}

func TestLintChartByPath(t *testing.T) {
	home1 := test.CreateTmpHome()
	home2 := test.CreateTmpHome()

	chartName := "goodChart"
	action.Create(chartName, home1)

	output := test.CaptureOutput(func() {
		Cli().Run([]string{"helmc", "--home", home2, "lint", util.WorkspaceChartDirectory(home1, chartName)})
	})

	test.ExpectContains(t, output, fmt.Sprintf("Chart [%s] has passed all necessary checks", chartName))
}

func TestLintAll(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	missingReadmeChart := "missingReadme"

	action.Create(missingReadmeChart, tmpHome)
	os.Remove(util.WorkspaceChartDirectory(tmpHome, missingReadmeChart, "README.md"))

	action.Create("goodChart", tmpHome)

	output := test.CaptureOutput(func() {
		Cli().Run([]string{"helmc", "--home", tmpHome, "lint", "--all"})
	})

	test.ExpectMatches(t, output, "A README file was not found.*"+missingReadmeChart)
	test.ExpectContains(t, output, "Chart [goodChart] has passed all necessary checks")
	test.ExpectContains(t, output, "Chart [missingReadme] failed some checks")
}
