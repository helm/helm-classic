package action

import (
	"io/ioutil"
	"testing"

	"github.com/helm/helm-classic/test"
	"github.com/helm/helm-classic/util"
)

func TestGenerate(t *testing.T) {
	ch := "generate"
	homedir := test.CreateTmpHome()
	test.FakeUpdate(homedir)
	Fetch(ch, ch, homedir)

	Generate(ch, homedir, []string{"ignore"}, true)

	// Now we should be able to load and read the `pod.yaml` file.
	path := util.WorkspaceChartDirectory(homedir, "generate/manifests/pod.yaml")
	d, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	pod := string(d)
	test.ExpectContains(t, pod, "image: ozo")
	test.ExpectContains(t, pod, "name: www-server")
}
