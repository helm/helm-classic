# Helm Plugins

Helm supports a basic plugin mechanism not unlike the one found in Git
or similar CLI projects. This feature is still considered experimental.

The basic model: When `helm` receives a subcommand that it does not know
(e.g. `helm foo`), it will look on `$PATH` for an executable named
`helm-foo`. If it finds one, it will set several environment variables
and then execute the named command, returning the results directly to
STDOUT and STDERR. Any flags passed after `foo` are passed on to the
`helm-foo` command. (Flags before foo, such as `helm -v foo`, are
interpreted by Helm. They may influence the environment, but are not
passed on.)

The plugin `plugins/sec/helm-sec` provides an example of how plugins can
be built.
