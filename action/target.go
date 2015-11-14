package action

import (
	"github.com/deis/helm/kubectl"
	"github.com/deis/helm/log"
)

// Target displays information about the cluster
func Target() {
	info, err := kubectl.ClusterInfo()
	if err != nil {
		log.Err("Could not connect to target: %s", err)
	}

	log.Msg(info)
}
