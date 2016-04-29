package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

const publishDescription = `This copies a chart from the workdir into the cache. Doing so
is the first stage of contributing a chart upstream.
`

var publishCmd = cli.Command{
	Name:        "publish",
	Usage:       "Publish a named chart to the git checkout.",
	Description: publishDescription,
	ArgsUsage:   "[chart-name]",
	Action: func(c *cli.Context) {
		minArgs(c, 1, "publish")
		action.Publish(c.Args()[0], home(c), c.String("repo"), c.Bool("force"))
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Force publish over an existing chart.",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "Publish to a specific chart repository.",
		},
	},
}
