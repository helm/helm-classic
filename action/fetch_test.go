package action

import (
	"testing"

	"github.com/helm/helm-classic/test"
	"github.com/helm/helm-classic/util"
)

func TestFetch(t *testing.T) {
	tmpHome := test.CreateTmpHome()
	test.FakeUpdate(tmpHome)
	chartName := "kitchensink"

	actual := test.CaptureOutput(func() {
		Fetch(chartName, "", tmpHome)
	})

	workspacePath := util.WorkspaceChartDirectory(tmpHome, chartName)
	test.ExpectContains(t, actual, "Fetched chart into workspace "+workspacePath)
}
