package action

func Install(chart, home string) {
	Fetch(chart, home)
	Info("kubectl create -f %s.yaml", chart)
}
