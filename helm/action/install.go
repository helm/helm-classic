package action

func Install(chart, home, namespace string) {
	Fetch(chart, chart, home, namespace)
	Info("kubectl create -f %s.yaml", chart)
}
