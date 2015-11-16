package kubectl

func ClusterInfo() (string, error) {
	return newKubectlCommand([]string{"cluster-info"}...).exec()
}
