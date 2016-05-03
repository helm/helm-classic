# Helm Classic Example Plugin

This example plugic highlights a few of the features of the Helm Classic plugin system.

## Usage

In this directory, build the example plugin:

```
$ go build -o helm-example helm-example.go
```

Helm Classic will search the path for plugins. So assuming you have Helm
installed, you can test your plugin like this:

```
$ PATH=$PATH:. helmc example -a foo -b bar baz
Args are: [helm-example -a foo -b bar baz]
Helm Classic home is: /Users/mattbutcher/.helmc
Helm Classic command is: example
Helm Classic default repo is: charts
```

The output of `helm-example` shows the contents of the arguments and
environment variables that Helm Classic passes to plugins.
