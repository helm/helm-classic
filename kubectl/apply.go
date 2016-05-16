package kubectl

// Apply uploads a chart to Kubernetes
func (r RealRunner) Apply(stdin []byte, ns string) ([]byte, error) {
	args := []string{"apply", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}

	cmd := command(args...)
	assignStdin(cmd, stdin)

	return cmd.CombinedOutput()
}

// Apply returns the commands to kubectl
func (r PrintRunner) Apply(stdin []byte, ns string) ([]byte, error) {
	args := []string{"apply", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}

	cmd := command(args...)
	assignStdin(cmd, stdin)

	return []byte(cmd.String()), nil
}
