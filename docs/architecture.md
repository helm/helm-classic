# Helm, Version 1

**TL;DR:** Homebrew for Kubernetes, with packages called "Charts"

This document outlines the architecture for a Kubernetes package management tool.

Our work suggests that great benefits could accrue if the community shares (and iterates on) a set of "best practices" Kubernetes packages. Examples of this might be:

- A Postgres database
- A cluster of Memcached servers
- A Varnish caching proxy

By providing common service definitions, lifecycle handling, health checking, and configuration, a Kubernetes package can serve the purposes of many Kubernetes users, and serve as a foundation for many other users.

## What Is a Kubernetes Package?

For our purposes, an "application" is any pre-packaged set of manifest files that can be deployed into Kubernetes, and that implement a particular functional unit of service (e.g. a web application, a database, or a message queue).

When it comes to deploying an application into Kubernetes there are always one or more Kubernetes manifest files. Often, these files (particularly pod and replication controller manifests) will reference one or more container images. These are typically pulled from external registries like GCR, Quay.io, and Docker Hub.

A Kubernetes package is thus:

- A collection of one or more  Kubernetes manifest files...
- That can be installed together...
- Together with descriptive metadata...
- And possibly some configuration scripts

In this document, we refer to a Kubernetes package as a Chart. The package format is discussed later in the document.

## A Package Archive

Right now, the Kubernetes community largely hand-crafts their manifest files. There is no accepted distribution schema, and the only commonly referenced set of files are the ones in the core Kubernetes GitHub repository.

A package-centric GitHub repository will provide a central place to store Charts. Users will fetch information about available packages from this service.

## The Workflows

This section describes a few different workflows

### User Workflow (Simple)

Use case: User is looking for a standalone etcd instance (one pod, no clustering). User wants it immediately installed on Kubernetes.

```
$ helm update
---> Fetching updates...
---> Done
$ helm search etcd
	etcd-standalone: etcd-standalone is a single Etcd instance (no cluster)
	etcd-ha-cluster: an HA Etcd cluster (3+ pods)
	etcd-discovery:  a single-pod Etcd service discovery server
$ helm info etcd-standalone
	Description: etcd-standalone is a single Etcd instance (no cluster)
	Version: 2.2.0-beta3
	Website: https://github.com/coreos/etcd
	Built: Oct. 21, 2015
	Provides: etcd-rc.yaml etcd-service.yaml
$ helm install etcd-standalone
---> Downloading etcd-standalone-2.2.0-beta3
---> Cached files into $KPM_HOME/etcd-standalone-2.2.0-beta3/
---> Running kubectl create -f ...
---> Done
```

### User Workflow (Advanced)

Use Case: User wants to get an NGINX formula, modify it, and then load the modified version into a running Kubernetes cluster.

```
$ helm update
---> Fetching updates...
---> Done
$ helm get nginx-standard
---> Downloading nginx-standard
---> Cached files into $HELM_HOME/nginx-standard/
$ vi $HELM_HOME/nginx-standard/nginx-rc.yaml
$ helm install nginx-standard
---> Running kubectl create -f ...
---> Done
```

TODO: If we back this to a strict Git workflow, we need a way to provide the user space for customized formulas.

## Helm Repository Design

This section describes the design of the Helm package repository.

The goals for designing an archive backend are as follows:

- Supports authentication and authorization
- Provides simple methods for storing, retrieving, searching, and browsing
- Intuitive to a software developer
- Highly available
- Understandable API
- Provides workflow for submitting, reviewing, and filing bugs

### GitHub Backend (Homebrew Model)

In this model, the only backend is a single GitHub repository. This repository contains all supported packages inside individual directories.

```
(root)
  |
  |- charts/
       |
       |- mychart/
       |
       |- postgres/
       |
       |- ...
```

The API to this service is essentially the git protocol.

Contributing packages to the archive is done by submitting GitHub Pull Requests.

## A Command Line Client

The command line client is the primary way that users interact with the Helm repository. The client is used for getting and installing packages.

The command line will support (at least) the following commands:

- `list`: List currently installed Charts
- `search PATTERN`: Search the available Charts, list those that match the `PATTERN`
- `get CHART`: Retrieve the Chart and put it in a local working directory
	* This will copy the Chart resource to a place where the user can modify it
- `build CHART`: Process the Chart
- `install CHART`: Get the Chart (if necessary) and then install it into Kubernetes, attaching to the k8s server configured for the current user via Kubectl.
- `info KUBE`: Print information about a Kube
- `update`: Refresh local metadata about available Kubes. On Homebrew model, this will re-pull the repo.
- `help`: Print help.

TODO: We would like to add a `test` command. How would this work?

The client will locally track the following three things:

- The client's own configuration
- The packages available from the remote service (a local cache)
- The packages that the client has "downloaded"

The tree looks something like this:

```
- $HELM_HOME
      |
      |- cache/       # Where `helm update` data goes
      |
      |- workdir/     # Working directory, where `helm get` copies go
      |
      |- config.yaml  # configuration info
```

The default location for `$HELM_HOME` will be `$HOME/.helm`. This reflects the fact that package management by this method is scoped to the developer, not to the system. However, we will make this flexible because we have the following target use cases in mind:

- Individual Helm developer/user
- CI/CD system
- Workdir shared by `git` repo among several developers
- Dockerized Helm that runs in a container

### Client Config

Clients will track minimal configuration about remote hosts, local state, and possibly security.

## A Web Interface

The initial web interface will focus exclusively on education:

- Introduction to Helm
- Instructions for installing
- Instructions for package submission

In the future, we will support the following additional features:

- Package search: Search for packages and show the details
- Test results: View the test results for any package

## CI/CD Service

To maintain a high bar for the submitted projects, a CI/CD service is run on each submitted project. The initial CI/CD implementation does the following:

- Runs the install script on an isolated container, and ensures that it exits with no errors
- Validates the manifest files
- Runs the package on its own in a Kubernetes cluster
- Verifies that the package can successfully enter the "Running" state (e.g. passes initial health checks)

## The Package Format

A package will be composed of the following parts:

- A metadata file
- One or more Kubernetes manifest files, stored in a `manifests/` directory

### The `Chart.yaml` File

A Manifest file is a YAML-formatted file describing a package's purpose, version, and authority information. It is called `Chart.yaml`.

The manifest file contains the following information:

- name: The name of the package
- version: A SemVer 2 version number, no leading `v`
- home: The URL to a relevant project page, git repo, or contact person
- description: A single line description
- details: A multi-line description
- maintainers: A list of name and URL/email address combinations for the maintainer(s)
- dependencies: A list of other services on which this depends. Dependencies contain names and version ranges. See _Dependecies_ below.

Example:

```yaml
name: ponycorn
home: http://github.com/technosophos/ponycorn
version: 1.0.3
description: The ponycorn server
maintainers:
  - Matt Butcher <mbutcher@deis.com>
dependencies:
  - name: nginx
    version: 1.9 < x <= 1.10
details:
  This package provides the ponycorn server, which requires an existing nginx
  server to function correctly.
```

#### Dependency Resolution

If dependencies are declared, `helm` will...

- Check to see if the named dependency is locally present
- If it is, check to see if the version is within the supplied parameters

If either check fails, `helm` will _emit a warning_, but will not prevent the installation. Under normal operation, it will prompt to continue. If the `--force` flag is set, it will simply continue.

Example:

```
$ helm install ponycorn
!!-> WARNING: Dependency nginx (1.9 < x <= 1.10) does not seem to be satisfied.
!!-> This package may not function correctly.
Continue? yes
---> Running `kubectl create -f ...
```

### Manifests

All Kubernetes manifests are stored in a `manifests/` directory. They must be in either YAML or JSON, and they must conform to the Kubernetes definitions.

Any pure YAML or JSON file must end with its appropriate suffix (`.yaml` or `.json`)

## Non-Features / Deferred Features

The following items have already been considered and will not be part of Version 1:

- Homebrew tap/keg system
- pre-install scripts to (for example) perform secrets generation
- post-install exec scripts (run `kubectl exec -it` on a pod)
