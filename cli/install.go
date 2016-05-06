package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
	"github.com/helm/helm-classic/kubectl"
)

const installDescription = `If the given 'chart-name' is present in your workspace, it
will be uploaded into Kubernetes. If no chart named 'chart-name' is found in
your workspace, Helm Classic will look for a chart with that name, install it into the
workspace, and then immediately upload it to Kubernetes.

When multiple charts are specified, Helm Classic will attempt to install all of them,
following the resolution process described above.
`

var installCmd = cli.Command{
	Name:        "install",
	Usage:       "Install a named package into Kubernetes.",
	Description: installDescription,
	ArgsUsage:   "[chart-name...]",
	Action:      install,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "namespace, n",
			Value: "",
			Usage: "The Kubernetes destination namespace.",
		},
		cli.BoolFlag{
			Name:  "force, aye-aye",
			Usage: "Perform install even if dependencies are unsatisfied.",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Fetch the chart, but only display the underlying kubectl commands.",
		},
		cli.BoolFlag{
			Name:  "generate,g",
			Usage: "Run the generator before installing.",
		},
		cli.StringSliceFlag{
			Name:  "exclude,x",
			Usage: "Files or directories to exclude from the generator (if -g is set).",
		},
	},
}

func install(c *cli.Context) {
	minArgs(c, 1, "install")
	h := home(c)
	force := c.Bool("force")

	client := kubectl.Client
	if c.Bool("dry-run") {
		client = kubectl.PrintRunner{}
	}

	for _, chart := range c.Args() {
		action.Install(chart, h, c.String("namespace"), force, c.Bool("generate"), c.StringSlice("exclude"), client)
	}
}
