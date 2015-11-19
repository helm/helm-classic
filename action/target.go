package action

import (
	"github.com/helm/helm/log"
)

// Target displays information about the cluster
func Target() {
	out, err := Kubectl.ClusterInfo()
	if err != nil {
		log.Err(err.Error())
	}
	log.Msg(string(out))
}
