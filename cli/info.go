package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

var infoCmd = cli.Command{
	Name:      "info",
	Usage:     "Print information about a Chart.",
	ArgsUsage: "[string]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "Print using a Go template",
		},
	},
	Action: func(c *cli.Context) {
		minArgs(c, 1, "info")
		action.Info(c.Args()[0], home(c), c.String("format"))
	},
}
