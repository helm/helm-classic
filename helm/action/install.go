package action

func Install(chart, home string) {
	Fetch(chart, chart, home)
	Info("kubectl create -f %s.yaml", chart)
}
