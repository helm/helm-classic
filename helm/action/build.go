package action

import (
	"github.com/deis/helm/helm/log"
)

func Build(chart, homedir string) {
	log.Info("Rebuilding %s", chart)
}
