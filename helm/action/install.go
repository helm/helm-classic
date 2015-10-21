package action

import (
	"github.com/deis/helm/helm/log"
)

func Install(chart, home, namespace string) {
	Fetch(chart, chart, home)
	log.Info("kubectl --namespace %q create -f %s.yaml", namespace, chart)
}
