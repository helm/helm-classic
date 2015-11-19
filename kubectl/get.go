package kubectl

import (
	"os/exec"
)

// Get returns Kubernetes resources
func (r RealRunner) Get(stdin []byte, ns string, dryRun bool) ([]byte, error) {
	args := []string{"get", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}
	return exec.Command("kubectl", args...).CombinedOutput()
}

func getCmd(args ...string) *exec.Cmd {
	return exec.Command("kubectl", args...)
}
