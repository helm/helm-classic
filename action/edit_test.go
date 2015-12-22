package action

import (
	"os"
	"path"
	"testing"

	"github.com/helm/helm/test"
)

func TestEdit(t *testing.T) {
	editor := os.Getenv("EDITOR")
	os.Setenv("EDITOR", "echo")
	defer os.Setenv("EDITOR", editor)

	tmpHome := test.CreateTmpHome()
	defer os.RemoveAll(tmpHome)
	test.FakeUpdate(tmpHome)

	Fetch("redis", "", tmpHome, false)

	expected := path.Join(tmpHome, "workspace/charts/redis")
	actual := test.CaptureOutput(func() {
		Edit("redis", tmpHome)
	})

	test.ExpectContains(t, actual, expected)
}
