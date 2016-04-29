package cli

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
	"github.com/helm/helm-classic/util"
)

var lintCmd = cli.Command{
	Name:      "lint",
	Usage:     "Validates given chart",
	ArgsUsage: "[chart-name]",
	Action:    lint,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all",
			Usage: "Check all available charts",
		},
	},
}

func lint(c *cli.Context) {
	home := home(c)

	all := c.Bool("all")

	if all {
		action.LintAll(home)
		return
	}

	minArgs(c, 1, "lint")

	a := c.Args()
	chartNameOrPath := a[0]

	fromHome := util.WorkspaceChartDirectory(home, chartNameOrPath)
	fromAbs := filepath.Clean(chartNameOrPath)

	_, err := os.Stat(fromAbs)

	if err == nil {
		action.Lint(fromAbs)
	} else {
		action.Lint(fromHome)
	}
}
