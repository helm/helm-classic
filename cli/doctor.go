package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/action"
)

const doctorDescription = `This will run a series of checks to ensure that your
experience with helmc is trouble-free.`

var doctorCmd = cli.Command{
	Name:        "doctor",
	Usage:       "Run a series of checks to surface possible problems",
	Description: doctorDescription,
	ArgsUsage:   "",
	Action: func(c *cli.Context) {
		action.Doctor(home(c))
	},
}
