package kubectl

import (
	"os/exec"
)

// Create uploads a chart to Kubernetes
func (r RealRunner) Create(stdin []byte, ns string) ([]byte, error) {
	args := []string{"create", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}

	cmd := exec.Command(Path, args...)
	assignStdin(cmd, stdin)

	return cmd.CombinedOutput()
}

// Create returns the commands to kubectl
func (r PrintRunner) Create(stdin []byte, ns string) ([]byte, error) {
	args := []string{"create", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}

	cmd := exec.Command(Path, args...)
	assignStdin(cmd, stdin)

	return []byte(commandToString(cmd)), nil
}
