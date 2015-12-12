package cli

import (
	"errors"

	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
	"github.com/helm/helm/log"
)

// version is the version of the app.
//
// This value is overwritten by the linker during build. The default version
// here is SemVer 2, but basically indicates that this was a one-off build
// and should not be trusted.
var version = "0.1.0"

const globalUsage = `The Kubernetes package manager

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

// Execute is the main entrypoint for the Helm CLI.
func Execute() {
	app := cli.NewApp()
	app.Name = "helm"
	app.Usage = globalUsage
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

	app.Commands = []cli.Command{
		createCmd,
		doctorCmd,
		editCmd,
		fetchCmd,
		homeCmd,
		infoCmd,
		installCmd,
		lintCmd,
		listCmd,
		publishCmd,
		removeCmd,
		repositoryCmd,
		searchCmd,
		targetCmd,
		uninstallCmd,
		updateCmd,
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

	app.Before = func(c *cli.Context) error {
		log.IsDebugging = c.Bool("debug")
		return nil
	}

	app.RunAndExitOnError()
}
