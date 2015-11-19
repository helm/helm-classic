package kubectl

import (
	"os/exec"
)

// Delete removes a chart from Kubernetes.
func (r RealRunner) Delete(name, ktype, ns string, dryRun bool) ([]byte, error) {

	args := []string{"delete", ktype, name}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}

	cmd := exec.Command(Path, args...)

	if dryRun {
		return []byte(commandToString(cmd)), nil
	}
	return cmd.CombinedOutput()
}
