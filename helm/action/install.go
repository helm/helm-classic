package action

func Install(chart, home, namespace string) {
	Fetch(chart, chart, home)
	Info("kubectl --namespace %q create -f %s.yaml", namespace, chart)
}
