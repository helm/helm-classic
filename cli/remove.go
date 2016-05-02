package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

var removeCmd = cli.Command{
	Name:      "remove",
	Aliases:   []string{"rm"},
	Usage:     "Remove one or more Charts from your working directory.",
	ArgsUsage: "[chart-name] [...]",
	Action:    remove,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force",
			Usage: "Remove Chart from working directory and leave packages installed.",
		},
	},
}

func remove(c *cli.Context) {
	minArgs(c, 1, "remove")
	h := home(c)
	force := c.Bool("force")

	a := c.Args()
	for _, chart := range a {
		action.Remove(chart, h, force)
	}

}
