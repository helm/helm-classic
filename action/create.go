package action

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/helm/helm-classic/chart"
	"github.com/helm/helm-classic/log"
	helm "github.com/helm/helm-classic/util"
)

// readmeSkel is the template for the README.md
const readmeSkel = `# {{.Name}}

Describe your chart here. Link to upstream repositories, Docker images or any
external documentation.

If your application requires any specific configuration like Secrets, you may
include that information here.
`

// manifestSkel is an example manifest for a new chart
const manifestSkel = `---
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

// Create a chart
//
// - chartName being created
// - homeDir is the helm home directory for the user
func Create(chartName, homeDir string) {
	chart := newSkelChartfile(chartName)
	createWithChart(chart, chartName, homeDir)
}

func createWithChart(chart *chart.Chartfile, chartName, homeDir string) {
	chartDir := helm.WorkspaceChartDirectory(homeDir, chartName)

	// create directories
	if err := os.MkdirAll(filepath.Join(chartDir, "manifests"), 0755); err != nil {
		log.Die("Could not create %q: %s", chartDir, err)
	}

	// create Chartfile.yaml
	if err := chart.Save(filepath.Join(chartDir, Chartfile)); err != nil {
		log.Die("Could not create Chart.yaml: err", err)
	}

	// create README.md
	if err := createReadme(chartDir, chart); err != nil {
		log.Die("Could not create README.md: err", err)
	}

	// create example-pod
	if err := createExampleManifest(chartDir); err != nil {
		log.Die("Could not create example manifest: err", err)
	}

	log.Info("Created chart in %s", chartDir)
}

// newSkelChartfile populates a Chartfile struct with example data
func newSkelChartfile(chartName string) *chart.Chartfile {
	return &chart.Chartfile{
		Name:        chartName,
		Home:        "http://example.com/your/project/home",
		Version:     "0.1.0",
		Description: "Provide a brief description of your application here.",
		Maintainers: []string{"Your Name <email@address>"},
		Details:     "This section allows you to provide additional details about your application.\nProvide any information that would be useful to users at a glance.",
	}
}

// createReadme populates readmeSkel and saves to the chart directory
func createReadme(chartDir string, c *chart.Chartfile) error {
	tmpl := template.Must(template.New("info").Parse(readmeSkel))

	readmeFile, err := os.Create(filepath.Join(chartDir, "README.md"))
	if err != nil {
		return err
	}

	return tmpl.Execute(readmeFile, c)
}

// createExampleManifest saves manifestSkel to the manifests directory
func createExampleManifest(chartDir string) error {
	return ioutil.WriteFile(filepath.Join(chartDir, "manifests/example-pod.yaml"), []byte(manifestSkel), 0644)
}
