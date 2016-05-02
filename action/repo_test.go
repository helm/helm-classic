package action

import (
	"testing"

	"github.com/helm/helm-classic/log"
	"github.com/helm/helm-classic/test"
)

func TestListRepos(t *testing.T) {
	log.IsDebugging = true

	homedir := test.CreateTmpHome()
	test.FakeUpdate(homedir)

	actual := test.CaptureOutput(func() {
		ListRepos(homedir)
	})

	test.ExpectContains(t, actual, "charts*\thttps://github.com/helm/charts")
}
