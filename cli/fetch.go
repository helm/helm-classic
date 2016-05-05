package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

const fetchDescription = `Copy a chart from the Chart repository to a local workspace.
From this point, the copied chart may be safely modified to your needs.

If an optional 'chart-name' is specified, the chart will be copied to a directory
of that name. For example, 'helmc fetch nginx www' will copy the the contents of
the 'nginx' chart into a directory named 'www' in your workspace.`

var fetchCmd = cli.Command{
	Name:        "fetch",
	Usage:       "Fetch a Chart to your working directory.",
	Description: fetchDescription,
	ArgsUsage:   "[chart] [chart-name]",
	Action:      fetch,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "namespace, n",
			Value: "default",
			Usage: "The Kubernetes destination namespace.",
		},
	},
}

func fetch(c *cli.Context) {
	home := home(c)
	minArgs(c, 1, "fetch")

	a := c.Args()
	chart := a[0]

	var lname string
	if len(a) == 2 {
		lname = a[1]
	}

	action.Fetch(chart, lname, home)
}
