package action

import (
	"fmt"
	"os/exec"
)

// Target displays information about the cluster
func Target() {
	CheckKubePrereqs()

	c, _ := exec.Command("kubectl", "cluster-info").Output()
	fmt.Println(string(c))
}
