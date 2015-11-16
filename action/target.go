package action

import (
	"fmt"
	"os/exec"

	"github.com/helm/helm/log"
)

// Target displays information about the cluster
func Target() {
	if _, err := exec.LookPath("kubectl"); err != nil {
		log.Die("Could not find 'kubectl' on $PATH: %s", err)
	}

	c, _ := exec.Command("kubectl", "cluster-info").Output()
	fmt.Println(string(c))
}
