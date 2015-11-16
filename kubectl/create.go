package kubectl

func Create(data []byte, ns string, dryRun bool) (string, error) {
	args := []string{"create", "-f", "-"}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}

	return newKubectlCommand(args...).withStdinData(string(data)).DryRun(dryRun).exec()
}
