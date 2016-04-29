package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

const listDescription = `This prints all of the packages that are currently installed in
the workspace. Packages are printed by the local name.
`

var listCmd = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List all fetched packages.",
	Description: listDescription,
	ArgsUsage:   "",
	Action: func(c *cli.Context) {
		action.List(home(c))
	},
}
