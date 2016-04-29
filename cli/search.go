package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

const searchDescription = `This provides a simple interface for searching the chart cache
for charts matching a given pattern.

If no string is provided, or if the special string '*' is provided, this will
list all available charts.
`

var searchCmd = cli.Command{
	Name:        "search",
	Usage:       "Search for a package.",
	Description: searchDescription,
	ArgsUsage:   "[string]",
	Action:      search,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "regexp,r",
			Usage: "Use a regular expression instead of a substring match.",
		},
	},
}

func search(c *cli.Context) {
	term := ""
	if len(c.Args()) > 0 {
		term = c.Args()[0]
	}
	action.Search(term, home(c), c.Bool("regexp"))
}
