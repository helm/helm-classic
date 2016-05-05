package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

const updateDescription = `This will synchronize the local repository with the upstream GitHub project.
The local cached copy is stored in '~/.helmc/cache' or (if specified)
'$HELMC_HOME/cache'.

The first time 'helmc update' is run, the necessary directory structures are
created and then the Git repository is pulled in full.

Subsequent calls to 'helmc update' will simply synchronize the local cache
with the remote.`

// updateCmd represents the CLI command for fetching the latest version of all charts from Github.
var updateCmd = cli.Command{
	Name:        "update",
	Aliases:     []string{"up"},
	Usage:       "Get the latest version of all Charts from GitHub.",
	Description: updateDescription,
	ArgsUsage:   "",
	Action: func(c *cli.Context) {
		if !c.Bool("no-version-check") {
			action.CheckLatest(version)
		}
		action.Update(home(c))
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "no-version-check",
			Usage: "Disable Helm Classic's automatic check for newer versions of itself.",
		},
	},
}
