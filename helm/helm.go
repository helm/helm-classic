package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/deis/helm/helm/action"
	"github.com/deis/helm/helm/log"
)

const version = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "helm"
	app.Usage = "The Kubernetes package manager"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo",
			Value:  "https://github.com/deis/helm",
			Usage:  "The remote Git repository as an HTTP URL",
			EnvVar: "HELM_REPO_URL",
		},
		cli.StringFlag{
			Name:   "home",
			Value:  "$HOME/.helm",
			Usage:  "The location of your Helm files",
			EnvVar: "HELM_HOME",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable verbose debugging output",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "update",
			Usage:     "Get the latest version of all Charts from GitHub.",
			ArgsUsage: "",
			Action:    update,
		},
		{
			Name:      "fetch",
			Usage:     "Fetch a Chart to your working directory.",
			ArgsUsage: "[chart] [chart-name]",
			Action:    fetch,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace, n",
					Value: "default",
					Usage: "The Kubernetes destination namespace.",
				},
			},
		},
		{
			Name:      "build",
			Usage:     "(Re-)build a manifest from templates.",
			ArgsUsage: "[chart-name...]",
			Action:    build,
		},
		{
			Name:      "install",
			Usage:     "Install a named package into Kubernetes.",
			ArgsUsage: "[chart-name...]",
			Action:    install,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace, n",
					Value: "",
					Usage: "The Kubernetes destination namespace.",
				},
			},
		},
		{
			Name:      "uninstall",
			Usage:     "Uninstall a named package from Kubernetes.",
			ArgsUsage: "[chart-name...]",
			Action:    uninstall,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace, n",
					Value: "",
					Usage: "The Kubernetes destination namespace.",
				},
			},
		},
		{
			Name:      "create",
			Usage:     "Create a chart in the local workspace",
			ArgsUsage: "[chart-name]",
			Action:    create,
		},
		{
			Name:      "edit",
			Usage:     "Edit a named chart in the local workspace",
			ArgsUsage: "[chart-name]",
			Action:    edit,
		},
		{
			Name:      "list",
			Usage:     "List all fetched packages",
			ArgsUsage: "",
			Action:    list,
		},
		{
			Name:      "search",
			Usage:     "Search for a package",
			ArgsUsage: "[string]",
			Action:    search,
		},
		{
			Name:      "info",
			Usage:     "Print information about a Chart",
			ArgsUsage: "[string]",
			Action:    info,
		},
		{
			Name:      "target",
			Usage:     "Displays information about cluster",
			ArgsUsage: "",
			Action:    target,
		},
	}

	app.Run(os.Args)
}

func home(c *cli.Context) string {
	return os.ExpandEnv(c.GlobalString("home"))
}

func repo(c *cli.Context) string {
	return os.ExpandEnv(c.GlobalString("repo"))
}

func update(c *cli.Context) {
	action.Update(repo(c), home(c))
}

func list(c *cli.Context) {
	action.List(home(c))
}

func fetch(c *cli.Context) {
	home := home(c)
	minArgs(c, 1, "fetch")

	a := c.Args()
	chart := a[0]

	var lname string
	if len(a) == 2 {
		lname = a[1]
	}

	action.Fetch(chart, lname, home)
}

func build(c *cli.Context) {
	home := home(c)

	for _, chart := range c.Args() {
		action.Build(chart, home)
	}
}

func install(c *cli.Context) {
	minArgs(c, 1, "install")
	for _, chart := range c.Args() {
		action.Install(chart, home(c), c.String("namespace"))
	}
}

func uninstall(c *cli.Context) {
	minArgs(c, 1, "uninstall")
	for _, chart := range c.Args() {
		action.Uninstall(chart, home(c), c.String("namespace"))
	}
}

func create(c *cli.Context) {
	minArgs(c, 1, "create")
	action.Create(c.Args()[0], home(c))
}

func edit(c *cli.Context) {
	minArgs(c, 1, "edit")
	action.Edit(c.Args()[0], home(c))
}

func search(c *cli.Context) {
	minArgs(c, 1, "search")
	action.Search(c.Args()[0], home(c))
}

func info(c *cli.Context) {
	minArgs(c, 1, "info")

	action.Info(c.Args()[0], home(c))
}

func target(c *cli.Context) {
	action.Target()
}

// minArgs checks to see if the right number of args are passed.
//
// If not, it prints an error and quits.
func minArgs(c *cli.Context, i int, name string) {
	if len(c.Args()) < i {
		m := "arguments"
		if i == 1 {
			m = "argument"
		}
		log.Err("Expected %d %s", i, m)
		cli.ShowCommandHelp(c, name)
		log.Die("")
	}
}
