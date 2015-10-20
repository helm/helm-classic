package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/deis/helm/helm/action"
	pretty "github.com/deis/pkg/prettyprint"
	//"github.com/technosophos/k8splace/model"
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
	}

	app.Commands = []cli.Command{
		{
			Name:   "update",
			Usage:  "Get the latest version of all Charts from GitHub.",
			Action: update,
		},
		{
			Name:   "fetch",
			Usage:  "Fetch a named package.",
			Action: fetch,
		},
		{
			Name:   "build",
			Usage:  "(Re-)build a manifest from templates.",
			Action: build,
		},
		{
			Name:   "install",
			Usage:  "Install a named package into Kubernetes.",
			Action: install,
		},
		{
			Name:   "list",
			Usage:  "List all fetched packages",
			Action: list,
		},
		{
			Name:   "search",
			Usage:  "Search for a package",
			Action: search,
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
	//h := NewClient(c.GlobalString("repo"))
	action.Info("Not implemented yet")
}

func fetch(c *cli.Context) {
	home := home(c)

	a := c.Args()

	if len(a) == 0 {
		action.Die("Fetch requires at least a Chart name")
	}

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
	for _, chart := range c.Args() {
		action.Install(chart, home(c))
	}
}

func search(c *cli.Context) {
	action.Info("Not implemented yet")
}

func info(msg string, args ...interface{}) {
	t := fmt.Sprintf(msg, args...)
	m := "{{.Yellow}}[INFO]{{.Default}} " + t
	fmt.Println(pretty.Colorize(m))
}

func ftw(msg string, args ...interface{}) {
	t := fmt.Sprintf(msg, args...)
	m := "{{.Green}}[YAY!]{{.Default}} " + t
	fmt.Println(pretty.Colorize(m))
}

func die(err error) {
	m := "{{.Red}}[BOO!]{{.Default}} " + err.Error()
	fmt.Println(pretty.Colorize(m))
	os.Exit(1)
}
