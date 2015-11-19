package kubectl

import (
	"os/exec"
)

// Get returns Kubernetes resources
func (r RealRunner) Get(stdin []byte, ns string) ([]byte, error) {
	args := []string{"get", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}
	cmd := exec.Command(Path, args...)
	assignStdin(cmd, stdin)

	return cmd.CombinedOutput()
}

// Get returns the commands to kubectl
func (r PrintRunner) Get(stdin []byte, ns string) ([]byte, error) {
	args := []string{"get", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}
	cmd := exec.Command(Path, args...)
	assignStdin(cmd, stdin)

	return []byte(commandToString(cmd)), nil
}
