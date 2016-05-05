package cli

import (
	"errors"

	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
	"github.com/helm/helm-classic/log"
)

// version is the version of the app.
//
// This value is overwritten by the linker during build. The default version
// here is SemVer 2, but basically indicates that this was a one-off build
// and should not be trusted.
var version = "0.1.0"

const globalUsage = `Helm Classic - A Kubernetes package manager

To begin working with Helm Classic, run the 'helmc update' command:

$ helmc update

This will download all of the necessary data. Common actions from this point
include:

- helmc help COMMAND: see help for a specific command
- helmc search: search for charts
- helmc fetch: make a local working copy of a chart
- helmc install: upload the chart to Kubernetes

For more information on Helm Classic, go to http://helm.sh.

ENVIRONMENT:
$HELMC_HOME:     Set an alternative location for Helm files. By default, these
				are stored in ~/.helmc

`

// Cli is the main entrypoint for the Helm Classic CLI.
func Cli() *cli.App {
	app := cli.NewApp()
	app.Name = "helmc"
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
			Value:  "$HOME/.helmc",
			Usage:  "The location of your Helm Classic files",
			EnvVar: "HELMC_HOME",
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
		generateCmd,
		tplCmd,
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

	return app
}
