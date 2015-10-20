package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
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
			Name:   "homedir",
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
			Action: list,
		},
	}

	app.Run(os.Args)
}

func update(c *cli.Context) {
}

func list(c *cli.Context) {
	//h := NewClient(c.GlobalString("repo"))
}

func fetch(c *cli.Context) {
}

func install(c *cli.Context) {
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

func ensureHome(c *cli.Context) string {
	wd := c.GlobalString("homedir")
	wd = os.ExpandEnv(wd)
	if _, err := os.Stat(wd); err != nil {
		info("Attempting to create dir %q", wd)
		if err := os.MkdirAll(wd, 0755); err != nil {
			die(err)
		}
		ftw("Created")
	}
	return wd
}

func die(err error) {
	m := "{{.Red}}[BOO!]{{.Default}} " + err.Error()
	fmt.Println(pretty.Colorize(m))
	os.Exit(1)
}
