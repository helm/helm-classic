package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/deis/helm/action"
	"github.com/deis/helm/log"
)

const version = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "helm"
	app.Usage = `The Kubernetes package manager

To begin working with Helm, run the 'helm update' command:

	$ helm update

This will download all of the necessary data. Common actions from this point
include:

	- helm help COMMAND: see help for a specific command
	- helm search: search for charts
	- helm fetch: make a local working copy of a chart
	- helm install: upload the chart to Kubernetes

For more information on Helm, go to http://helm.sh.

ENVIRONMENT:
   $HELM_HOME:     Set an alternative location for Helm files. By default, these
                   are stored in ~/.helm
   $HELM_REPO_URL: Set an alternative upstream chart repository.

`
	app.Version = version
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo",
			Value:  "https://github.com/deis/charts",
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

	app.CommandNotFound = func(c *cli.Context, command string) {
		log.Err("No matching command '%s'", command)
		cli.ShowAppHelp(c)
		log.Die("")
	}

	app.Commands = []cli.Command{
		{
			Name:  "update",
			Usage: "Get the latest version of all Charts from GitHub.",
			Description: `This will synchronize the local repository with the upstream GitHub project.
The local cached copy is stored in '~/.helm/cache' or (if specified)
'$HELM_HOME/cache'.

The first time 'helm update' is run, the necessary directory structures are
created and then the Git repository is pulled in full.

Subsequent calls to 'helm update' will simply synchronize the local cache
with the remote.`,
			ArgsUsage: "",
			Action: func(c *cli.Context) {
				action.Update(repo(c), home(c))
			},
		},
		{
			Name:  "fetch",
			Usage: "Fetch a Chart to your working directory.",
			Description: `Copy a chart from the Chart repository to a local workspace.
From this point, the copied chart may be safely modified to your needs.

If an optional 'chart-name' is specified, the chart will be copied to a directory
of that name. For example, 'helm fetch nginx www' will copy the the contents of
the 'nginx' chart into a directory named 'www' in your workspace.
`,
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
			Name:      "remove",
			Aliases:   []string{"rm"},
			Usage:     "Removes a Chart from your working directory.",
			ArgsUsage: "[chart-name]",
			Action:    remove,
		},
		{
			Name:  "install",
			Usage: "Install a named package into Kubernetes.",
			Description: `If the given 'chart-name' is present in your workspace, it
will be uploaded into Kubernetes. If no chart named 'chart-name' is found in
your workspace, Helm will look for a chart with that name, install it into the
workspace, and then immediately upload it to Kubernetes.

When multiple charts are specified, Helm will attempt to install all of them,
following the resolution process described above.

As a special case, if the flag --chart-path is specified, Helm will look for a
Chart.yaml file and manifests/ directory at that path, and will install that
chart if found. In this case, you may not specify multiple charts at once.
`,
			ArgsUsage: "[chart-name...]",
			Action:    install,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "chart-path,p",
					Value: "",
					Usage: "An alternate path to fetch a chart. If specified, Helm will ignore the cache.",
				},
				cli.StringFlag{
					Name:  "namespace, n",
					Value: "",
					Usage: "The Kubernetes destination namespace.",
				},
				cli.BoolFlag{
					Name:  "force, aye-aye",
					Usage: "Perform install even if dependencies are unsatisfied.",
				},
			},
		},
		{
			Name:  "uninstall",
			Usage: "Uninstall a named package from Kubernetes.",
			Description: `For each supplied 'chart-name', this will connect to Kubernetes
and remove all of the manifests specified.

This will not alter the charts in your workspace.
`,
			ArgsUsage: "[chart-name...]",
			Action: func(c *cli.Context) {
				minArgs(c, 1, "uninstall")
				for _, chart := range c.Args() {
					action.Uninstall(chart, home(c), c.String("namespace"))
				}
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace, n",
					Value: "",
					Usage: "The Kubernetes destination namespace.",
				},
			},
		},
		{
			Name:  "create",
			Usage: "Create a chart in the local workspace.",
			Description: `This will scaffold a new chart named 'chart-name' in your
local workdir. To edit the resulting chart, you may edit the files directly or
use the 'helm edit' command.
`,
			ArgsUsage: "[chart-name]",
			Action: func(c *cli.Context) {
				minArgs(c, 1, "create")
				action.Create(c.Args()[0], home(c))
			},
		},
		{
			Name:  "edit",
			Usage: "Edit a named chart in the local workspace.",
			Description: `Existing charts in the workspace can be edited using this command.
'helm edit' will open all of the chart files in a single editor (as specified
by the $EDITOR environment variable).
`,
			ArgsUsage: "[chart-name]",
			Action: func(c *cli.Context) {
				minArgs(c, 1, "edit")
				action.Edit(c.Args()[0], home(c))
			},
		},
		{
			Name:  "publish",
			Usage: "Publish a named chart to the git checkout.",
			Description: `This copies a chart from the workdir into the cache. Doing so
is the first stage of contributing a chart upstream.
`,
			ArgsUsage: "[chart-name]",
			Action: func(c *cli.Context) {
				action.Publish(c.Args()[0], home(c), c.Bool("force"))
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Force publish over an existing chart.",
				},
			},
		},
		{
			Name:  "list",
			Usage: "List all fetched packages.",
			Description: `This prints all of the packages that are currently installed in
the workspace. Packages are printed by the local name.
`,
			ArgsUsage: "",
			Action: func(c *cli.Context) {
				action.List(home(c))
			},
		},
		{
			Name:  "search",
			Usage: "Search for a package.",
			Description: `This provides a simple interface for searching the chart cache
for charts matching a given pattern.

If no string is provided, or if the special string '*' is provided, this will
list all available charts.
`,
			ArgsUsage: "[string]",
			Action:    search,
		},
		{
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
		},
		{
			Name:      "target",
			Usage:     "Displays information about cluster.",
			ArgsUsage: "",
			Action: func(c *cli.Context) {
				action.Target()
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		log.IsDebugging = c.Bool("debug")
		return nil
	}

	app.Run(os.Args)
}

// home runs the --home flag through os.ExpandEnv.
func home(c *cli.Context) string {
	return os.ExpandEnv(c.GlobalString("home"))
}

// repo runs the --repo flag through os.ExpandEnv.
func repo(c *cli.Context) string {
	return os.ExpandEnv(c.GlobalString("repo"))
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

func remove(c *cli.Context) {
	home := home(c)
	minArgs(c, 1, "remove")

	a := c.Args()
	chart := a[0]

	action.Remove(chart, home)
}

func install(c *cli.Context) {
	minArgs(c, 1, "install")
	h := home(c)
	force := c.Bool("force")

	// If chart-path is specified, we do an alternative install.
	//
	// This version will only install one chart at a time, since the
	// chart-path can only point to one chart.
	if alt := c.String("chart-path"); alt != "" {
		action.AltInstall(c.Args()[0], alt, h, c.String("namespace"), force)
		return
	}

	for _, chart := range c.Args() {
		action.Install(chart, h, c.String("namespace"), force)
	}
}

func search(c *cli.Context) {
	term := ""
	if len(c.Args()) > 0 {
		term = c.Args()[0]
	}
	action.Search(term, home(c))
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
