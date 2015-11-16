# Helm Example Plugin

This plugin highlights a few of the features of the Helm plugin system.

## Usage

In this directory, build the example plugin:

```
$ go build -o helm-example helm-example.go
```

Helm will search the path for plugins. So assuming you have Helm
installed, you can test your plugin like this:

```
$ PATH=$PATH:. helm example -a foo -b bar baz
Args are: [helm-example -a foo -b bar baz]
Helm home is: /Users/mattbutcher/Code/helm_home
Helm command is: example
Helm default repo is: charts
```

The output of `helm-example` shows the contents of the arguments and
environment variables that Helm passes to plugins.
