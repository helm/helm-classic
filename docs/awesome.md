# The Helm Classic Guide to Writing Awesome Charts

A Helm Classic Chart provides a recipe for installing and running a containerized application inside of Kubernetes. This guide explains how to write an outstanding Chart.

## Creating a New Chart

You can create a new chart using the `helmc create` command. This will put the chart in your workspace, which is the perfect place for trying it out. You may choose to use the `helmc edit` command to edit your chart, or you may be more comfortable editing directly with your favorite editor.

## A Brief Anatomy of a Chart

Inside of your newly created chart, there are three pieces:

- The `Chart.yaml` file contains descriptive data about a chart.
- The `manifests/` directory contains one or more Kubernetes manifest files.
- A `README.md` file, which is a Markdown-formatted file that explains how to configure your chart.

This guide walks you through the process of creating those parts.

## Naming Your Chart

A good name succinctly expresses what your chart provides. If your chart installs an application in the "normal" way or is a basic configuration of a package, you should name it with the name of the project. For example, a singe MongoDB chart may simple be called `mongo`.

If your chart displays some more advanced features, or creates a cluster of servers, you should give it a more specific name: `redis-cluster` creates and configures a cluster of Redis servers.

Please take a look at the existing charts before naming your own.

_A Note on example charts:_ As Helm Classic was getting started, we created a few example charts, whose purpose was to illustrate how to write a chart. We are trying to only add new _example_ charts when they illustrate something new and helpful for chart developers.

## The Chart.yaml File

A simple Chart.yaml file should look like this:

```
name: NAME OF CHART
home: URL TO PROJECT HOME
source:
  - URL TO IMAGE SOURCE
version: VERSION
description: SHORT DESCRIPTION
dependencies:
  - name: CHART
    from: GIT URL
    version: SEMVER FILTER
maintainers:
  - YOUR NAME <YOUR EMAIL>
details:
  ONE PARAGRAPH DESCRIPTION
```

- name: This should be the same name as your chart directory.
- home: The URL to the project where the container image came from. This is to assist people in finding out how their image was built.
- source: A set of URLs to the source of container images
- version: A SemVer 2 version string.
- description: A short (several word) description of this chart.
- dependencies: A set of name/filter pairs
	- name: The name of the chart this chart depends upon
  - repo: A Git-flavor URL pointing to the repository of origin (e.g. `https://github.com/foo/bar.git`)
	- version: A filter indicating what version of the chart is required. Example: `~1.2` (greater than or equal to 1.2.0, and less than 1.3.0)
- maintainers: A set of maintainer names, together with an email
- details: A single paragraph describing the chart

Except for `dependencies`, all fields are required.

### Dependency Resolution

When a chart is fetched or installed, Helm Classic will perform dependency
resolution and alert the user if the chart dependencies are not
satisfied.

When just a `name` is provided, Helm Classic will verify that a chart by that
name exists in the same Git repo, and that it is fetched into the
workspace.

When `name` and `repo` is provided, Helm Classic will verify that a chart by
that name exists in the given repo, and is fetched into the current
workspace.

When `version` is added (in either case), Helm Classic will additionally verify
that the chart in the workspace has a version within the bounds of the
specified version. Remember that the `version` section can us version
ranges, fuzzy versions, and [so on](https://github.com/Masterminds/semver#hyphen-range-comparisons).

## The README.md File

The README file performs one important function: It tells the user how to use your chart. It is automatically displayed when a chart is installed.

Your README.md file _may_:

- Tell the user how to configure the chart (if necessary)
- How to work with the chart once it's in Kubernetes (if this is not self-evident)
- Explain any environment variables, secrets, annotations, or special labels
- Provide URLs to external sources that are helpful

Your README.md file _should not_:

- Add a license constraint. All charts in `github.com/helm/charts` are Apache v2. (If you're using your own repository, you may add a license if you wish.)
- Describe the project in detail

In a nutshell, the README is intended as help text for your user: It gets the user from fetching your chart to using your chart. You may choose to start with [a README template](https://gist.github.com/elsonrodriguez/6e9afea8c0afbaf9cf59) and go from there.

## Manifest Files

All Kubernetes manifest files should be in YAML format. We know some people prefer JSON, but we've decided for a number of reasons to adopt YAML as the single format for Helm Classic manifests.

A good simple manifest typically includes:

- A replication controller (`foo-rc.yaml`)
- A service definition for each backend service this chart relies upon (`foo-db-svc.yaml`)
- A secrets file for each secret the RC uses (`foo-password-sec.yaml`)

### Keeper Manifests

To support upgrading between versions of a chart, Helm Classic allows individual manifests to be "keepers." These manifests get special treatment:

- `helmc uninstall` skips them while printing a warning
- `helmc install` applies changes rather than always creating a new manifest

Marking a manifest as a keeper is accomplished by adding a `helm-keep: "true"` annotation:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: deis-router
  namespace: deis
  labels:
    heritage: deis
  annotations:
    helm-keep: "true"
```

This mechanism allows essential pieces of infrastructure to remain in place while other components are uninstalled and reinstalled. For example, a chart might mark its `Namespace` and any externally visible `Service` manifests with `helm-keep`, to ensure that DNS entries aren't invalidated by destroying the `LoadBalancer` or `NodePort` ingress.

The `helm-keep` annotation is respected by Helm Classic version 0.8.0 and later.

### Labels

All Helm Classic charts should have an `app` label and a `heritage: helm` label in their metadata sections. These provide a base-level consistency across all Helm Classic charts. (`heritage: helm` makes it easy to search a Kubernetes cluster for all components installed via Helm Classic.)

We suggest using the following labels where appropriate:

- `app`: The name of the app that this manifest is part of.
- `tier`
- `provider`: Indicates that this chart provides a specific backend service (e.g. MariaDB or Percona might include `provider: mysql` since they provide MySQL compatible service)

### Environment Variables

We recommend documenting these in the README. Something like this is great:

```
- DB_PASSWORD: The password used to connect to the MySQL database
```

### Pods vs. Replication Controllers

From the Kubernetes documentation:

> [W]e recommend that you use a replication controller even if your application requires only a single pod

Because replication controllers come with better lifecycle management, we suggest using RCs instead of using pods. This will give your users all of the features of a pod, but with added assurances. That said, when you use an RC to run a service that _cannot be scaled_, you should indicate this to your users either by adding the `scalable: "false"` label to the RC, or by saying so directly in the README.

### Secrets

We recommend providing default secrets whenever a secret is needed, and also explaining the secret in the README.

### Jobs, DaemonSets, and other experimental features

When using a definition from the experimental set, please note this in the `Chart.yaml` description with the notation `(experimental)`.

Example:

```yaml
name: sweep
description: Job to sweep the cluster once. (experimental)
# ...
details:
  This sweeps out the dusty corners of the cluster. It uses the experimental Jobs type.
```

## Testing

The Helm Classic chart reviewers expect that you have tested your chart, both for compatibility with Helm, and also with Kubernetes. We run some basic tests on your charts on submission, but these perform only rudimentary checks. So please test before submitting.

## Publishing a Pull Request

The suggested workflow for publishing a pull request goes like this:

1. From GitHub, clone the `github.com/helm/charts` repository
2. Add your fork to Helm Classic: `helm repo add technosophos git@github.com:technosophos/charts.git`
3. Copy your chart from your workspace to your new charts repo (see `helmc publish`)
4. Commit and push to your charts repo
5. Using GitHub, file a PR (Pull Request) against the `helm/charts` repository
6. Follow along in the `helm/charts` issue queue

## Review

The core Helm Classic charts team reviews all charts to see if they comply with our chart best practices. **All charts must be reviewed and marked LGTM by two members of the Helm Classic charts team.** Once that is done, your chart will be merged into the repository.

We absolutely _love_ contributions, so don't fret about this part of the process. When we ask for additional changes, it's only because we (like you) want the Helm Classic community to have the best experience possible.

See you in the issue queue!
