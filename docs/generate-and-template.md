# Helm Classic Generate and Template

Introduced in Helm Classic 0.3.0, Helm Classic has the ability to embed and run
generators that can perform arbitrary modifications on charts. The
intent of this feature is to provide chart developers with the ability
to modify charts using strategies like parameterization.

Along with the generation feature, Helm Classic includes a templating command
that simply provides Go template support.

This document describes how to use this pair of tools, and how to extend
generation to use other tools.

## Templates

Helm Classic includes a simple templating solution based on Go templates.
Running `helmc template example.yaml` will compile and execute
the template `example.yaml`.

As a simple example, here is a basic template:

```
apiVersion: v1
kind: Namespace
metadata:
  name: {{lower "Foo"}}
  labels:
    heritage: deis
```

The above is a Kubernetes namespace YAML file with one template
directive: `{{lower "Foo"}}` will lowercase the string "Foo" to "foo" and
insert it into the template. Here is the result of running `helmc
template example.yaml`:

```
apiVersion: v1
kind: Namespace
metadata:
  name: foo
  labels:
    heritage: deis
```

### Parameterizing Templates

Templates become more powerful when we parameterize them. For example,
we can modify the above template to say "If there is a parameter named
`.Namespace`, use that. Otherwise set the default namespace to 'foobar'":

```
apiVersion: v1
kind: Namespace
metadata:
  name: {{default "foobar" .Namespace}}
  labels:
    heritage: deis
```

Once again, if we render that template, we'll see the same results:

```
apiVersion: v1
kind: Namespace
metadata:
  name: foobar
  labels:
```

### Template Functions

Helm Classic's template tool includes an array of built-in functions. A few
particularly useful ones are:

- `default <default> <other>`: This function allows you to set a default value. We saw it
  used above.
- `b64enc <string>`: This base-64 encodes a string. Useful for secrets.
- `randAlphaNum <int>`: This generates a random alphanumeric string of
  length INT
- `env <string>`: This retrieves the value of an environment variable.

But there are 40+ functions available. These functions are supplied by
the built-in Go `text/template` package and the [Sprig template function
library](https://github.com/Masterminds/sprig).

### The Values File

To pass parameters into the `helmc template` rendering process, we need
to supply a file of values. Above we used the `.Namespace` placeholder,
but no value was ever supplied to it, and so it is always empty.

We can now supply a value in a file called `values.toml`.

```
Namespace = "goldie"
```

Helm Classic templates support three file formats for values files, and it uses
the file extensions to determine which type of data is in a value file.

- [TOML](https://github.com/toml-lang/toml) is a clear and concise
  configuration file format (`.toml`).
- YAML is the same format Helm Classic uses (`.yaml` or `.yml`).
- JSON is the Javascript-like syntax (`.json`).

We prefer TOML because it is clear and concise, and because the
different format makes it easy to visually distinguish values files from
manifest files. But you can use whichever you like.

Once we have a values file, we can tell `helmc template` about it:

```
$ helmc template -d values.toml example.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: goldie
  labels:
    heritage: deis
```

Note that `name` is now `goldie` (the value in `values.toml`) instead of
the default `foobar`.

Finally, by default `helmc template` writes to STDOUT. You can redirect
output to a file with the `-o` or `--out` flags:

```
$ helmc template -o out.yaml -d values.toml example.yaml
```

This final form is the one used most frequently by generators.

## Generators

Helm Classic provides a command called `helmc generate`. This function operates
on charts within your workspace (e.g. charts that you have fetched or
installed).

The generator's job is to _scan the files inside of a chart looking for
`helm:generate` headers_.

Files that begin with the appropriate `helm:generate` header will
trigger a generation run.

### The `helm:generate` Header

The `helm:generate` header must be the first line of a file (any file),
and must exactly follow one of the three formats described below:

```
#helm:generate CMD [ARGS]
//helm:generate CMD [ARGS]
/*helm:generate CMD [ARGS]*/
```

The generate header is space-sensitive and case-sensitive. Lines MUST
begin with one of the above comment sequences, and MUST have a lowercase
`helm:generate` string.

`CMD` is the command that Helm Classic will run, and the optional `[ARGS]` are
arguments passed to the `CMD`.

For example, we could embed a simple `sed` script inside of a generate
header:

```
#helm:generate sed -i -e s|ubuntu-debootstrap|fluffy-bunny| my/pod.yaml
```

It is important to note that generate commands are _not run in a shell_.
However, environment variables are expanded.

Along with any existing environment variables, the following variables
are specially defined:

- `HELM_HOME`: The Helm Classic home directory
    - Although the value of the `HELMC_HOME` environment variable, if defined, may play a role in defining the specially expanded `HELM_HOME` variable, the two are not equivalent. To produce charts that remain compatible with the _original_ Helm tool, which has now become Helm Classic, only the `HELM_HOME` variable should be used within `helm:generate` headers.
- `HELM_DEFAULT_REPO`: The repository alias for the default repository.
- `HELM_GENERATE_FILE`: The present file's name
- `HELM_GENERATE_DIR`: The absolute path to the chart directory of the present chart

These are available both from the invocation of the `CMD`, and from
inside any generator command itself.

As an example, consider a generation header like this, inside of a file
named `example.js`:

```
//helm:generate echo $HELM_GENERATE_FILE
```

This will be expanded to the command `echo /path/to/example.json`.

### Writing A Custom Generator

A generator is any tool that is executable within your environment. When
`helmc generate` runs, it looks for the generator in two places:

- First, if the generator is an absolute path (any path that begins with
  `/`), it will execute it on that path.
- Otherwise it will search for `CMD` on `$PATH`.

To execute a script relative to the chart's path, you can use
the `$HELM_GENERATE_DIR` variable:

```
/*helm:generate $HELM_GENERATE_DIR/myscript.sh */
```

Note that the above variable is the absolute path to the chart
directory, not to the file that `helmc generate` found.

The following guidelines should be kept in mind when writing Helm Classic
generator plugins:

- All of the named environment variables are available to the plugin.
- The plugin does _not_ execute in a shell. Thus you cannot expect that
  implementors will be able to easily redirect STDOUT or STDERR using
  shell redirects.
  - It is advised, therefore, that you not write generated output to
   STDOUT, but provide a flag for writing to a file.
  - User-facing output like progress meters _should_ be written to
  STDOUT. `helmc generate` will pass this information on to the user
  during execution.
  - Errors should be written to STDERR so that `helmc generate` can
   capture them and redirect them to the user.
- Plugins should return the exit code 0 on success. Any other exit code
  will be treated as an error and will stop execution of further
  generators.
- No recovery action is performed by `helmc generate`. So if your script
  dirties the environment, it must also clean the environment.

In addition to these, it is important to remember that the plugins you
write may not be used in isolation. Other plugins may operate on the
same chart, even during the course of the same `helmc generate` run.

It is considered a best practice to _not_ embed generator headers inside
of `Chart.yaml` or inside of `manifests/` files. Conversely, it is
considered good practice to encourage plugin users to store files in
separate subdirectories inside of a chart. For example, if your tool is
called `secgenerate`, you might suggest that any charts that use this
tool place their tool-specific files within `$CHART/secgenerate/`.

## Combining Generate and Template

The `helmc template` command is an example of a generator plugin. It is
designed to be invoked within a generator. Here is an example of using
the two in conjunction:

The [Deis Namespace]() chart contains a `tpl/` directory with a file
named `namespace.yaml`. This file contains the following text:

```yaml
#helm:generate helm tpl -d tpl/values.toml -o manifests/namespace.yaml $HELM_GENERATE_FILE
apiVersion: v1
kind: Namespace
metadata:
  name: {{default "deis" .Namespace}}
  labels:
    heritage: deis
```

__IMPORTANT NOTE:__ To produce charts that remain compatible with the _original_ Helm tool, which has now become Helm Classic, the `helmc` binary should not be invoked _directly_ from within a `helm:generator` header. Using `helm` instead will allow older versions of the tool to execute your generators successfully. Meanwhile, Helm Classic intelligently invokes `helmc` wherever `helm` is invoked within a `helm:generator` header.

To run the generators, we first fetch the chart, and then run the
generator:

```
$ helmc fetch deis/namespace
---> Fetched chart into workspace /Users/mattbutcher/.helmc/workspace/charts/namespace
---> Done
$  helmc generate namespace
---> Ran 2 generators.
```

At this point, if we look inside of the generated `namespace.yaml` file,
we will see the result:

```
$ cat $(helmc home)/workspace/charts/namespace/manifests/namespace.yaml
#helm:generate helm tpl -d tpl/values.toml -o manifests/namespace.yaml $HELM_GENERATE_FILE
apiVersion: v1
kind: Namespace
metadata:
  name: deis
  labels:
    heritage: deis
```

To modify the value of the `name`, we merely need to create a
`values.toml` file in the workspace's `charts/namespace/tpl/` directory
and we will be able to substitute.

## Conclusion

Generators are a powerful way of extending Helm Classic's generative
capabilities. While we have provided the template generator, the system
can be extended to your needs.
