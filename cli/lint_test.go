package cli

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/helm/helm/action"
	"github.com/helm/helm/test"
	"github.com/helm/helm/util"
)

func TestLintFlags(t *testing.T) {
	actualFlags := make([]string, len(lintCmd.Flags))

	for i, flag := range lintCmd.Flags {
		actualFlags[i] = flag.String()
	}

	expectedFlags := []string{"--all	Check all available charts"}

	if !reflect.DeepEqual(actualFlags, expectedFlags) {
		t.Fatalf("Expected: %v, Actual: %v", expectedFlags, actualFlags)
	}
}

func TestLintAllNone(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	output := test.CaptureOutput(func() {
		Cli().Run([]string{"helm", "--home", tmpHome, "lint", "--all"})
	})

	test.ExpectContains(t, output, fmt.Sprintf("Could not find any charts in \"%s", tmpHome))
}

func TestLintSingleCli(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)

	chartName := "goodChart"
	action.Create(chartName, tmpHome)

	output := test.CaptureOutput(func() {
		Cli().Run([]string{"helm", "--home", tmpHome, "lint", chartName})
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
		Cli().Run([]string{"helm", "--home", tmpHome, "lint", "--all"})
	})

	test.ExpectMatches(t, output, "A README file was not found.*"+missingReadmeChart)
	test.ExpectContains(t, output, "Chart [goodChart] has passed all necessary checks")
	test.ExpectContains(t, output, "Chart [missingReadme] failed some checks")
}
