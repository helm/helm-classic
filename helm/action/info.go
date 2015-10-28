package action

import (
	"path/filepath"
	"text/template"

	"github.com/deis/helm/helm/log"
	"github.com/deis/helm/helm/model"
)

const defaultInfoFormat = `Name: {{.Name}}
Home: {{.Home}}
Version: {{.Version}}
Description: {{.Description}}
Details: {{.Details}}
`

// Print information about a chart
//
// - chartName to display
// - homeDir is the helm home directory for the user
// - format is a optional Go template
func Info(chartName, homedir, format string) {
	chartPath := filepath.Join(homedir, CacheChartPath, chartName, "Chart.yaml")

	if format == "" {
		format = defaultInfoFormat
	}

	chart, err := model.LoadChartfile(chartPath)
	if err != nil {
		log.Die("Could not find chart %s \nError %s", chartName, err.Error())
	}

	log.Info(chartName)

	tmpl, err := template.New("info").Parse(format)
	if err != nil {
		log.Die("%s", err)
	}

	if err = tmpl.Execute(log.Stdout, chart); err != nil {
		log.Die("%s", err)
	}
}
