# Using the Workspace

Your local Helm Classic home directory looks like this:

```
$HELMC_HOME
├── cache               # The cache of all existing chart repositories
│   ├── charts          # The cache of the helm/charts repository
│   │   ├── .git        # Each cached repository is a git repository
│   │   ├── ...
│   │   ├── mysql
│   │   ├── redis
│   │   └── ...
│   ├── deis            # An example of another cached chart repository
│   │   ├── .git        # Each cached repository is a git repository
│   │   ├── ...
│   │   ├── workflow-dev
│   │   └── ...
│   └── ...
└── workspace
    └── charts          # Charts that have been fetched from the cache
        ├── redis
        ├── workflow-dev
        └── ...
```

In this document, we focus on the `workspace` directory. We suggest some ways to make the most of your Workspace. But before we get to that, let's take a quick look at the `cache` directory.

## The Cache Directory

The `cache` contains local clones of remote chart repositories, such as [github.com/helm/charts](https://github.com/helm/charts).  Anytime you run a `helmc update`, Helm Classic will re-sync each of these clone to their respective upstreams.

When you are developing Charts to be contributed upstream, you will interact with the cloned repositories in the `cache` directory directly. But in normal day-to-day usage, you do not need to worry about it.

## The Workspace

"Infrastructure is code." One of the giant benefits of Kubernetes is that we can use a simple declarative model to describe microservice infrastructure. Helm Classic is designed to make it easier to manage your Kubernetes manifest files.

At its simplest, the workspace is the place where your Charts go. When you run the command `helmc fetch nginx`, it copies the `nginx` chart from the cache to your workspace. There, you can modify it to your heart's content, customizing the manifests for your particular needs.

The `helmc install` command also copies a chart into your workspace (only if you don't already have it there) so that even your quick builds are reproducible. And the mention of reproducibility leads us to the main topic at hand: _the workspace describes your infrastructure_.

### Dependencies, Naming, and Copying

When `helmc fetch` or `helmc install` fetches a chart into the workspace,
it automatically adds the `from:` section to a chart. For example,
running `helmc fetch alpine myalpine` will create a `Chart.yaml` in
`$HELMC_HOME/workspace/charts/myalpine`. That chart will look like this:

```
name: myalpine
from:
  name: alpine
  version: 0.1.1
  repo: https://github.com/deis/charts
# ...
```

The `from:` section contains information about the base chart that was
fetched into the workspace. This information is important because it is
used during dependency resolution.

- `name` is the name of the base chart
- `version` is the version of the base chart that was fetched
- `repo` is the Git repository from which the base chart was fetched

These fields match one-to-one with the fields that can be specified in
the `dependency` section of a chart.

## Best Practice for your Workspace

Most Helm Classic users spend at least a little bit of time experimenting. They run a few installs, edit a few charts, and see what they can do. But we hope that at some point users transition from experimentation to real-world usage.

Here's where the Helm Classic team has decided to be hands-off so that you can institute practices that are best for you.

### Scenario 1: The Workspace Team Repo

A small team uses several charts to deploy Kubernetes objects with Helm Classic. Some of these charts are unmodified, but many more are tailored specifically to the team's needs.

How should the team keep track of these?

Our suggestion is to create a source code repository (using a tool like Git or Mercurial) mapped to your workspace.

```
$ cd $HELMC_HOME/workspace
$ git init
$ git remote add...
$ git push -u origin master
```

Now the workspace is tracked by a version control system, and the team can begin sharing charts as configuration. This provides three distinct advantages:

1. This makes sharing charts among a team easy, convenient, and familiar.
2. Since your local copy of the charts is stored, if you accidentally overwrite or break a chart, you can restore a previous version.
3. Your CI/CD pipeline can take advantage of version control hooks, using charts to deploy to testing, staging, or production.

With all of these advantages, why didn't we (the Helm Classic authors) just create this repository for you? We thought long and hard about it, but we decided that teams are better equipped to decide how to set up this tooling for their own environment.

### Scenario 2: The Chart Developer

Helm Classic comes with tooling to make chart development easy. If you want to create a new chart, it is as easy as this:

```
$ helmc create mychart
```

This will scaffold out a new chart named `mychart` will be created in `$HELMC_HOME/workspace/charts/mychart`. This directory will have a basic `Chart.yaml` and a `manifests` directory.

You can easily get started editing with the `helmc edit mychart` command.

Charts are edited in the workspace, giving you a place to edit, test, and tune without having to dirty your `cache` copy (a lesson we learned from corrupting our own Homebrew repos).

Once your chart is ready for submission, you can set up your `cache` directory and submit upstream:

1. Fork the `github.com/helm/charts` repository
2. Set up your local cache to track the fork: `cd $HELMC_HOME/cache && git remote add ...`
3. Run `helmc publish` to push your chart into your cache
4. Commit, push, and issue your pull request.

## Conclusion

The workspace is designated as a location that Helm Classic does not explicitly manage. Helm Classic does not manage it so that you can pick the tools best for you and your team, and also tailor the environment to your use case.
