package action

import (
	"text/template"

	"github.com/helm/helm-classic/chart"
	"github.com/helm/helm-classic/log"
	helm "github.com/helm/helm-classic/util"
)

const defaultInfoFormat = `Name: {{.Name}}
Home: {{.Home}}
Version: {{.Version}}
Description: {{.Description}}
Details: {{.Details}}
`

// Info prints information about a chart.
//
// - chartName to display
// - homeDir is the helm home directory for the user
// - format is a optional Go template
func Info(chartName, homedir, format string) {
	r := mustConfig(homedir).Repos
	table, chartLocal := r.RepoChart(chartName)
	chartPath := helm.CacheDirectory(homedir, table, chartLocal, Chartfile)

	if format == "" {
		format = defaultInfoFormat
	}

	chart, err := chart.LoadChartfile(chartPath)
	if err != nil {
		log.Die("Could not find chart %s: %s", chartName, err.Error())
	}

	tmpl, err := template.New("info").Parse(format)
	if err != nil {
		log.Die("%s", err)
	}

	if err = tmpl.Execute(log.Stdout, chart); err != nil {
		log.Die("%s", err)
	}
}
