package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

const createDescription = `This will scaffold a new chart named 'chart-name' in your
local workdir. To edit the resulting chart, you may edit the files directly or
use the 'helmc edit' command.`

var createCmd = cli.Command{
	Name:        "create",
	Usage:       "Create a chart in the local workspace.",
	Description: createDescription,
	ArgsUsage:   "[chart-name]",
	Action: func(c *cli.Context) {
		minArgs(c, 1, "create")
		action.Create(c.Args()[0], home(c))
	},
}
