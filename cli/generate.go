package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
)

var generateCmd = cli.Command{
	Name:      "generate",
	Usage:     "Run the generator over the given chart.",
	ArgsUsage: "[chart-name]",
	Action: func(c *cli.Context) {
		home := home(c)
		minArgs(c, 1, "generate")

		a := c.Args()
		chart := a[0]
		action.Generate(chart, home)
	},
}
