package action

import (
	"github.com/helm/helm-classic/kubectl"
	"github.com/helm/helm-classic/log"
)

// Target displays information about the cluster
func Target(client kubectl.Runner) {
	out, err := client.ClusterInfo()
	if err != nil {
		log.Err(err.Error())
	}
	log.Msg(string(out))
}
