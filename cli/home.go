package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm/log"
)

var homeCmd = cli.Command{
	Name:      "home",
	Usage:     "Displays the location of the Helm home.",
	ArgsUsage: "",
	Action: func(c *cli.Context) {
		log.Msg(home(c))
	},
}
