package action

import (
	"io/ioutil"
	"testing"

	"github.com/helm/helm-classic/test"
	"github.com/helm/helm-classic/util"
)

func TestCreate(t *testing.T) {
	tmpHome := test.CreateTmpHome()

	Create("mychart", tmpHome)

	// assert chartfile
	chartfile, err := ioutil.ReadFile(util.WorkspaceChartDirectory(tmpHome, "mychart/Chart.yaml"))
	if err != nil {
		t.Errorf("Could not read chartfile: %s", err)
	}
	actualChartfile := string(chartfile)
	expectedChartfile := `name: mychart
home: http://example.com/your/project/home
version: 0.1.0
description: Provide a brief description of your application here.
maintainers:
- Your Name <email@address>
details: |-
  This section allows you to provide additional details about your application.
  Provide any information that would be useful to users at a glance.
`
	test.ExpectEquals(t, actualChartfile, expectedChartfile)

	// asset readme
	readme, err := ioutil.ReadFile(util.WorkspaceChartDirectory(tmpHome, "mychart/README.md"))
	if err != nil {
		t.Errorf("Could not read README.md: %s", err)
	}
	actualReadme := string(readme)
	expectedReadme := `# mychart

Describe your chart here. Link to upstream repositories, Docker images or any
external documentation.

If your application requires any specific configuration like Secrets, you may
include that information here.
`
	test.ExpectEquals(t, expectedReadme, actualReadme)

	// assert example manifest
	manifest, err := ioutil.ReadFile(util.WorkspaceChartDirectory(tmpHome, "mychart/manifests/example-pod.yaml"))
	if err != nil {
		t.Errorf("Could not read manifest: %s", err)
	}
	actualManifest := string(manifest)
	expectedManifest := `---
apiVersion: v1
kind: Pod
metadata:
  name: example-pod
  labels:
    heritage: helm
spec:
  restartPolicy: Never
  containers:
  - name: example
    image: "alpine:3.2"
    command: ["/bin/sleep","9000"]
`
	test.ExpectEquals(t, actualManifest, expectedManifest)
}
