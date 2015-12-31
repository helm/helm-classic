package cli

import (
	"github.com/codegangsta/cli"
	"github.com/helm/helm/action"
)

var generateDesc = `Read a chart and generate manifests based on generator tools.

'helm generate' reads all of the files in a chart, searching for a generation
header. If it finds the header, it will execute the corresponding command.

The header is in the form '#helm:generate CMD [ARGS]'. 'CMD' can be any command
that Helm finds on $PATH. Optional 'ARGS' are arguments that will be passed on
to the command. Generator commands can begin with any of the following sequences:
'#helm:generate', '//helm:generate', or '/*helm:generate'.

For example, to embed a generate instruction in a YAML file, one may do the
following:

	#helm:generate helm tpl mytemplate.yaml

If CMD is an absolute path, Helm will attempt to execute it even if it is not
on $PATH. Combined with the $HELM_GENERATE_DIR environment variable, charts can
include their own local scripts:

	#helm:generate $HELM_GENERATE_DIR/myscript.sh

Since '#' is a comment character in YAML, the YAML parser will ignore the
generator line. But 'helm:generate' will read it as specifying that the following
command should be run:

	helm tpl mytemplate.yaml

While the 'helm tpl' command can easily be used in conjunction with the
'helm generate' command, you are not limited to just this tool. For example, one
could run a sed substitution just as easily:

	#helm:generate sed -i -e s|ubuntu-debootstrap|fluffy-bunny| my/pod.yaml

Note that 'helm generate' does not execute inside of a shell. However, it does
expand environment variables. The following variables are made available by the
Helm system:

- HELM_HOME: The Helm home directory
- HELM_DEFAULT_REPO: The repository alias for the default repository.
- HELM_GENERATE_FILE: The present file's name
- HELM_GENERATE_DIR: The absolute path to the chart directory of the present chart

By default, 'helm generate' will execute every generator that it finds in a
project. Generators can be mixed, with different files using different
generators. The order of generation is the order in which the directory contents
are listed.

The environment variables listed above are also available to generators.

For charts that contain multiple different generator template sets, you may
prevent generators from being run using the '--exclude' flag:

	$ helm generate --exclude=tpl --exclude=sed foo

The above will prevent the generator from traversing the 'foo' chart's 'tpl/'
or 'sed/' directories.
`

var generateCmd = cli.Command{
	Name:        "generate",
	Usage:       "Run the generator over the given chart.",
	ArgsUsage:   "[chart-name]",
	Description: generateDesc,
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "exclude,x",
			Usage: "Files or directories to exclude from this run, relative to the chart.",
		},
	},
	Action: func(c *cli.Context) {
		home := home(c)
		minArgs(c, 1, "generate")

		a := c.Args()
		chart := a[0]
		action.Generate(chart, home, c.StringSlice("exclude"))
	},
}
