package cli

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/log"
)

// home runs the --home flag through os.ExpandEnv.
func home(c *cli.Context) string {
	return os.ExpandEnv(c.GlobalString("home"))
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
