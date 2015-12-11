package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
)

const uninstallDescription = `For each supplied 'chart-name', this will connect to Kubernetes
	and remove all of the manifests specified.

	This will not alter the charts in your workspace.
	`

var uninstallCmd = cli.Command{
	Name:        "uninstall",
	Usage:       "Uninstall a named package from Kubernetes.",
	Description: uninstallDescription,
	ArgsUsage:   "[chart-name...]",
	Action: func(c *cli.Context) {
		minArgs(c, 1, "uninstall")
		for _, chart := range c.Args() {
			action.Uninstall(chart, home(c), c.String("namespace"), c.Bool("force"))
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "namespace, n",
			Value: "",
			Usage: "The Kubernetes destination namespace.",
		},
		cli.BoolFlag{
			Name:  "force, aye-aye, y",
			Usage: "Do not ask for confirmation.",
		},
	},
}
