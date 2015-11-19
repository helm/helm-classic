package kubectl

import (
	"os/exec"
)

// ClusterInfo returns Kubernetes cluster info
func (r RealRunner) ClusterInfo() ([]byte, error) {
	return exec.Command(Path, "cluster-info").CombinedOutput()
}

// ClusterInfo returns the commands to kubectl
func (r PrintRunner) ClusterInfo() ([]byte, error) {
	cmd := exec.Command(Path, "cluster-info")
	return []byte(commandToString(cmd)), nil
}
