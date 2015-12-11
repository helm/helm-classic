package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
)

const editDescription = `Existing charts in the workspace can be edited using this command.
	'helm edit' will open all of the chart files in a single editor (as specified
	by the $EDITOR environment variable).
	`

var editCmd = cli.Command{
	Name:        "edit",
	Usage:       "Edit a named chart in the local workspace.",
	Description: editDescription,
	ArgsUsage:   "[chart-name]",
	Action: func(c *cli.Context) {
		minArgs(c, 1, "edit")
		action.Edit(c.Args()[0], home(c))
	},
}
