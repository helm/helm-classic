package kubectl

import (
	"os/exec"
)

// Create uploads a chart to Kubernetes
func (r RealRunner) Create(stdin []byte, ns string, dryRun bool) ([]byte, error) {

	args := []string{"create", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}

	cmd := exec.Command(Path, args...)
	assignStdin(cmd, stdin)

	if dryRun {
		return []byte(commandToString(cmd)), nil
	}
	return cmd.CombinedOutput()
}
