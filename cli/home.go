package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/log"
)

var homeCmd = cli.Command{
	Name:      "home",
	Usage:     "Displays the location of the Helm Classic home.",
	ArgsUsage: "",
	Action: func(c *cli.Context) {
		log.Msg(home(c))
	},
}
