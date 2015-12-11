package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
)

const updateDescription = `This will synchronize the local repository with the upstream GitHub project.
The local cached copy is stored in '~/.helm/cache' or (if specified)
'$HELM_HOME/cache'.

The first time 'helm update' is run, the necessary directory structures are
created and then the Git repository is pulled in full.

Subsequent calls to 'helm update' will simply synchronize the local cache
with the remote.This will synchronize the local repository with the upstream GitHub project.
The local cached copy is stored in '~/.helm/cache' or (if specified)
'$HELM_HOME/cache'.

The first time 'helm update' is run, the necessary directory structures are
created and then the Git repository is pulled in full.

Subsequent calls to 'helm update' will simply synchronize the local cache
with the remote.`

// UpdateCmd represents the CLI command for fetching the latest version of all charts from Github.
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
			Usage: "Disable Helm's automatic check for newer versions of itself.",
		},
	},
}
