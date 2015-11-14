package kubectl

func Delete(name, ktype, ns string, dryRun bool) error {
	args := []string{"delete", ktype, name}

	if ns != "" {
		args = append([]string{"--namespace=" + ns}, args...)
	}

	_, err := newKubectlCommand(args...).DryRun(dryRun).exec()
	return err
}
