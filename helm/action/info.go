package action

import (
	"path/filepath"

	"github.com/deis/helm/helm/log"
	"github.com/deis/helm/helm/model"
)

func Info(chart, homedir string) {
	chartPath := filepath.Join(homedir, CacheChartPath, chart, "Chart.yaml")

	log.Info("%s", chartPath)

	chartModel, err := model.LoadChartfile(chartPath)
	if err != nil {
		log.Die("%s - UNKNOWN", chart)
	}

	log.Info("Chart: %s", chartModel.Name)
	log.Info("Description: %s", chartModel.Description)
	log.Info("Details: %s", chartModel.Details)
	log.Info("Version: %s", chartModel.Version)
	log.Info("Website: %s", chartModel.Home)
	log.Info("From: %s", chartPath)
	log.Info("Dependencies: %s", chartModel.Dependencies)
}
