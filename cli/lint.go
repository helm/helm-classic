package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
)

var lintCmd = cli.Command{
	Name:      "lint",
	Usage:     "Validates given chart",
	ArgsUsage: "[chart-name]",
	Action:    lint,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all",
			Usage: "Check all available charts",
		},
	},
}

func lint(c *cli.Context) {
	home := home(c)

	all := c.Bool("all")

	if all {
		action.LintAll(home)
	} else {
		minArgs(c, 1, "lint")

		a := c.Args()
		chart := a[0]

		action.Lint(chart, home)
	}
}
