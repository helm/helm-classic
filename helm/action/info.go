package action

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/deis/helm/helm/log"
	"github.com/deis/helm/helm/model"
)

func Info(chart, homedir string) {
	chartPath := filepath.Join(homedir, CacheChartPath, chart, "Chart.yaml")

	chartDescription, err := describeChart(chartPath)
	if err != nil {
		log.Die("Could not find chart %s", chart)
	}

	log.Info("%s", chartDescription)
}

func describeChart(chartPath string) (bytes.Buffer, error) {

	var output bytes.Buffer

	chartModel, err := model.LoadChartfile(chartPath)
	if err != nil {
		return output, err
	}

	w := bufio.NewWriter(&output)
	fmt.Fprintf(w, "Chart: %s\n", chartModel.Name)
	fmt.Fprintf(w, "Description: %s\n", chartModel.Description)
	fmt.Fprintf(w, "Details: %s\n", chartModel.Details)
	fmt.Fprintf(w, "Version: %s\n", chartModel.Version)
	fmt.Fprintf(w, "Website: %s\n", chartModel.Home)
	fmt.Fprintf(w, "From: %s\n", chartPath)
	fmt.Fprintf(w, "Dependencies: %s\n", chartModel.Dependencies)
	w.Flush()

	return output, nil
}
