package main

import (
	"errors"
	"os"

	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
	"github.com/helm/helm/log"
)

// version is the version of the app.
//
// This value is overwritten by the linker during build. The default version
// here is SemVer 2, but basically indicates that this was a one-off build
// and should not be trusted.
var version = "0.1.0-unstable"

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

`
	app.Version = version
	app.EnableBashCompletion = true
	app.After = func(c *cli.Context) error {
		if log.ErrorState {
			return errors.New("Exiting with errors")
		}

		return nil
	}

	app.Flags = []cli.Flag{
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
		if action.HasPlugin(command) {
			action.Plugin(home(c), command, c.Args())
			return
		}
		log.Err("No matching command '%s'", command)
		cli.ShowAppHelp(c)
		log.Die("")
	}

	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"up"},
			Usage:   "Get the latest version of all Charts from GitHub.",
			Description: `This will synchronize the local repository with the upstream GitHub project.
The local cached copy is stored in '~/.helm/cache' or (if specified)
'$HELM_HOME/cache'.

The first time 'helm update' is run, the necessary directory structures are
created and then the Git repository is pulled in full.

Subsequent calls to 'helm update' will simply synchronize the local cache
with the remote.`,
			ArgsUsage: "",
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
			Name:      "lint",
			Usage:     "Validates given chart",
			ArgsUsage: "[chart-name]",
			Action:    lint,
		},
		{
			Name:      "remove",
			Aliases:   []string{"rm"},
			Usage:     "Remove one or more Charts from your working directory.",
			ArgsUsage: "[chart-name] [...]",
			Action:    remove,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force",
					Usage: "Remove Chart from working directory and leave packages installed.",
				},
			},
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
					Name:  "namespace, n",
					Value: "",
					Usage: "The Kubernetes destination namespace.",
				},
				cli.StringFlag{
					Name:  "value, v",
					Value: "",
					Usage: "Specify a list of key-value pairs (eg. -v FOO=BAR,BAR=FOO) to set/override parameter values in Templates or Config Data.",
				},
				cli.StringFlag{
					Name:  "value-folder",
					Value: "chart-config",
					Usage: "Specify the folder to load and save any overridden parameter values for Templates or Config Data so that they can be reused across install/update commands.",
				},
				cli.BoolFlag{
					Name:  "force, aye-aye",
					Usage: "Perform install even if dependencies are unsatisfied.",
				},
				cli.BoolFlag{
					Name:  "dry-run",
					Usage: "Fetch the chart, but only display the underlying kubectl commands.",
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
					action.Uninstall(chart, home(c), c.String("namespace"), c.Bool("force"))
				}
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace, n",
					Value: "",
					Usage: "The Kubernetes destination namespace.",
				},
				cli.BoolFlag{
					Name:  "force, aye-aye, y",
					Usage: "Do not ask for confirmation.",
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
				minArgs(c, 1, "publish")
				action.Publish(c.Args()[0], home(c), c.String("repo"), c.Bool("force"))
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Force publish over an existing chart.",
				},
				cli.StringFlag{
					Name:  "repo",
					Usage: "Publish to a specific chart repository.",
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all fetched packages.",
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "regexp,r",
					Usage: "Use a regular expression instead of a substring match.",
				},
			},
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
		{
			Name:      "home",
			Usage:     "Displays the location of the Helm home.",
			ArgsUsage: "",
			Action: func(c *cli.Context) {
				log.Msg(home(c))
			},
		},
		{
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
		},
		{
			Name:  "doctor",
			Usage: "Run a series of checks to surface possible problems",
			Description: `This will run a series of checks to ensure that your
experience with helm is trouble-free.
`,
			ArgsUsage: "",
			Action: func(c *cli.Context) {
				action.Doctor(home(c))
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		log.IsDebugging = c.Bool("debug")
		return nil
	}

	app.RunAndExitOnError()
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

func lint(c *cli.Context) {
	home := home(c)
	minArgs(c, 1, "lint")

	a := c.Args()
	chart := a[0]

	action.Lint(chart, home)
}

func remove(c *cli.Context) {
	minArgs(c, 1, "remove")
	h := home(c)
	force := c.Bool("force")

	a := c.Args()
	for _, chart := range a {
		action.Remove(chart, h, force)
	}

}

func install(c *cli.Context) {
	minArgs(c, 1, "install")
	h := home(c)
	force := c.Bool("force")
	dryRun := c.Bool("dry-run")

	for _, chart := range c.Args() {
		action.Install(chart, h, c.String("namespace"), force, dryRun, c.String("value"), c.String("value-folder"))
	}
}

func search(c *cli.Context) {
	term := ""
	if len(c.Args()) > 0 {
		term = c.Args()[0]
	}
	action.Search(term, home(c), c.Bool("regexp"))
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
