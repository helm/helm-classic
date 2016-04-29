package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

var repositoryCmd = cli.Command{
	Name:    "repository",
	Aliases: []string{"repo"},
	Usage:   "Work with other Chart repositories.",
	Subcommands: []cli.Command{
		{
			Name:      "add",
			Usage:     "Add a remote chart repository.",
			ArgsUsage: "[name] [git url]",
			Action: func(c *cli.Context) {
				minArgs(c, 2, "add")
				a := c.Args()
				action.AddRepo(home(c), a[0], a[1])
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all remote chart repositories.",
			Action: func(c *cli.Context) {
				action.ListRepos(home(c))
			},
		},
		{
			Name:      "remove",
			Aliases:   []string{"rm"},
			Usage:     "Remove a remote chart repository.",
			ArgsUsage: "[name] [git url]",
			Action: func(c *cli.Context) {
				minArgs(c, 1, "remove")
				action.DeleteRepo(home(c), c.Args()[0])
			},
		},
	},
}
