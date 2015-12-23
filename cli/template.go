package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
)

const tplDescription = `Execute a template inside of a chart.

This command is not intended to be run directly (though it can be). Instead, it
is a helper for the generate command. Run 'helm help generate' for more.

'helm template' provides a default implementation of a templating feature for
Kubernetes manifests. Other more sophisticated methods can be plugged in using
the 'helm generate' system.

'helm template' uses Go's built-in text template system to provide template
substitution inside of a chart. In addition to the built-in template commands,
'helm template' supports all of the template functions provided by the Sprig
library (https://github.com/Masterminds/sprig).

If a values data file is provided, 'helm template' will use that as a source
for values. If none is specified, only default values will be used. Helm uses
simple extension scanning to determine the file type of the values data file.

- YAML: .yaml, .yml
- TOML: .toml
- JSON: .json

If an output file is specified, the results will be written to the output
file instead of STDOUT. Writing to the source template file is unsupported.
(In other words, don't set the source and output to the same file.)
`

// tplCmd is the command to handle templating.
// glide tpl -o dest.txt -d data.toml my_template.tpl
var tplCmd = cli.Command{
	Name:      "template",
	Aliases:   []string{"tpl"},
	Usage:     "Run a template command on a file.",
	ArgsUsage: "[file]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "out,o",
			Usage: "The destination file. If unset, results are written to STDOUT.",
		},
		cli.StringFlag{
			Name:  "values,d",
			Usage: "A file containing values to substitute into the template. TOML (.toml), JSON (.json), and YAML (.yaml, .yml) are supported.",
		},
	},
	Action: func(c *cli.Context) {
		minArgs(c, 1, "template")

		a := c.Args()
		filename := a[0]
		action.Template(c.String("out"), filename, c.String("values"))
	},
}
