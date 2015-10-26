# Using the Workspace

Your local Helm directory looks like this:

```
$HELM_HOME
├── cache
│   ├── .git            # Cache is a git repo
│   ├── charts          # The cache of all existing charts
|   |   ├── alpine
│   │   └── ...
│   ├── docs
│   └── helm            # The Helm client source
└── workspace
    └── charts
        ├── alpine
        ├── myserver
        └── ...
```

In this document, we focus on the `workspace` directory. We suggest some ways to make the most of your Workspace. But before we get to that, let's take a quick look at the `cache` directory.

## The Cache Directory

The `cache` directory is your local copy of the Helm repository. If you compare it to [github.com/deis/helm](github.com/deis/helm) they should look the same. And when you run a `helm update`, it will sync your local `cache` directory to the upstream.

When you are developing Charts to be contributed upstream, you will interact with the `cache` directory directly. But in normal day-to-day usage, you do not need to worry about it.

## The Workspace

"Infrastructure is code." One of the giant benefits of Kubernetes is that we can use a simple declarative model to describe microservice infrastructure. Helm is designed to make it easier to manage your Kubernetes manifest files.

At its simplest, the workspace is the place where your Charts go. When you run the command `helm fetch nginx`, it copies the `nginx` chart to your workspace. There, you can modify it to your heart's content, customizing the manifests for your particular needs.

The `helm install` command also copies a chart into your workspace (if you don't already have it there) so that even your quick builds are reproducible. And the mention of reproducibility leads us to the main topic at hand: _the workspace describes your infrastructure_.

## Best Practice for your Workspace

Most Helm users spend at least a little bit of time experimenting with Helm. They run a few installs, edit a few charts, and see what they can do. But we hope that at some point users transition from experimentation to real-world usage.

Here's where the Helm team has decided to be hands-off so that you can institute practices that are best for you.

### Scenario 1: The Workspace Team Repo

A small team uses several charts to deploy Kubernetes objects with Helm. Some of these charts are unmodified, but many more are tailored specifically to the team's needs.

How should the team keep track of these?

Our suggestion is to create a source code repository (using a tool like Git or Mercurial) mapped to your workspace.

```
$ cd $HELM_HOME/workspace
$ git init
$ git remote add...
$ git push -u origin master
```

Now the workspace is tracked by a version control system, and the team can begin sharing charts as configuration. This provides three distinct advantages:

1. This makes sharing charts among a team easy, convenient, and familiar.
2. Since your local copy of the charts is stored, if you accidentally overwrite or break a chart, you can restore a previous version.
3. Your CI/CD pipeline can take advantage of version control hooks, using charts to deploy to testing, staging, or production.

With all of these advantages, why didn't we (the Helm authors) just create this repository for you? We thought long and hard about it, but we decided that teams are better equipped to decide how to set up this tooling for their own environment.

### Scenario 2: The Chart Developer

Helm comes with tooling to make chart development easy. If you want to create a new chart, it is as easy as this:

```
$ helm create mychart
```

This will scaffold out a new chart named `mychart` will be created in `$HELM_HOME/workspace/charts/mychart`. This directory will have a basic `Chart.yaml` and a `manifests` directory.

You can easily get started editing with the `helm edit mychart` command.

Charts are edited in the workspace, giving you a place to edit, test, and tune without having to dirty your `cache` copy (a lesson we learned from corrupting our own Homebrew repos).

Once your chart is ready for submission, you can set up your `cache` directory and submit upstream:

1. Fork the `github.com/deis/helm` repository
2. Set up your local cache to track the fork: `cd $HELM_HOME/cache && git remote add ...`
3. Run `helm publish` to push your chart into your cache
4. Commit, push, and issue your pull request.

## Conclusion

The workspace is designated as a location that Helm does not explicitly manage. Helm does not manage it so that you can pick the tools best for you and your team, and also tailor the environment to your use case.