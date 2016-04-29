package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
	"github.com/helm/helm-classic/kubectl"
)

var targetCmd = cli.Command{
	Name:      "target",
	Usage:     "Displays information about cluster.",
	ArgsUsage: "",
	Action: func(c *cli.Context) {
		client := kubectl.Client
		if c.Bool("dry-run") {
			client = kubectl.PrintRunner{}
		}
		action.Target(client)
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Only display the underlying kubectl commands.",
		},
	},
}
