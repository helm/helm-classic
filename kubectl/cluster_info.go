package kubectl

import (
	"os/exec"
)

// ClusterInfo returns Kubernetes cluster info
func (r RealRunner) ClusterInfo() ([]byte, error) {
	return exec.Command(Path, "cluster-info").CombinedOutput()
}
