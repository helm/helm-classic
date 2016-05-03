# Helm Classic Plugins

Helm Classic supports a basic plugin mechanism not unlike the one found in Git
or similar CLI projects. This feature is still considered experimental.

The basic model: When `helmc` receives a subcommand that it does not know
(e.g. `helmc foo`), it will look on `$PATH` for an executable named
`helm-foo`. If it finds one, it will set several environment variables
and then execute the named command, returning the results directly to
STDOUT and STDERR. Any flags passed after `foo` are passed on to the
`helm-foo` command. (Flags before foo, such as `helmc -v foo`, are
interpreted by Helm Classic. They may influence the environment, but are not
passed on.)

The plugin `plugins/sec/helm-sec` provides an example of how plugins can
be built.
