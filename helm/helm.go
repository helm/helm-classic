package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/deis/helm/helm/action"
	pretty "github.com/deis/pkg/prettyprint"
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
					Value: "default",
					Usage: "The Kubernetes destination namespace.",
				},
			},
		},
		{
			Name:      "list",
			Usage:     "List all fetched packages",
			ArgsUsage: "",
			Action:    list,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace, n",
					Value: "default",
					Usage: "The Kubernetes destination namespace.",
				},
				cli.BoolFlag{
					Name:  "all-namespaces",
					Usage: "List all namespaces. Equivalent to -n '*'",
				},
			},
		},
		{
			Name:      "search",
			Usage:     "Search for a package",
			ArgsUsage: "[string]",
			Action:    search,
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
	if c.Bool("all-namespaces") {
		action.List(home(c), "*")
		return
	}
	action.List(home(c), c.String("namespace"))
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

	action.Fetch(chart, lname, home, c.String("namespace"))
}

func build(c *cli.Context) {
	home := home(c)

	for _, chart := range c.Args() {
		action.Build(chart, home)
	}
}

func install(c *cli.Context) {
	for _, chart := range c.Args() {
		action.Install(chart, home(c), c.String("namespace"))
	}
}

func search(c *cli.Context) {
	action.Search(c.Args()[0], home(c))
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
