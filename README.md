# Helm Classic - A Kubernetes Package Manager

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/156px-Warning.svg.png) | Helm Classic (the `helm/helm` repository) is **no longer actively developed** but will remain available until `kubernetes/helm` has stabilized.
|---|---|

Helm and [Deployment Manager](https://github.com/kubernetes/deployment-manager)
have recently joined forces to make deploying and managing software on
Kubernetes as simple as possible. The combined effort now lives in the Kubernetes GitHub organization at
[kubernetes/helm][k8s-helm].

For more information about the architecture of `kubernetes/helm` read the **[architecture documentation](https://github.com/kubernetes/helm/blob/master/docs/architecture.md)**.

## Overview

[![Build Status](https://travis-ci.org/helm/helm-classic.svg?branch=master)](https://travis-ci.org/helm/helm-classic) [![Go Report Card](https://goreportcard.com/badge/github.com/helm/helm-classic)](https://goreportcard.com/report/github.com/helm/helm-classic)

[Helm Classic](https://helm.sh) bootstraps your Kubernetes cluster with **Charts** that provide ready-to-use workloads like:

- A Redis cluster
- A Postgres database
- An HAProxy edge load balancer

A Chart is a unit of Kubernetes manifests that reflect best practices as determined by the Helm Classic community.  Helm Classic's [architecture](docs/architecture.md) is heavily influenced by [Homebrew](https://github.com/Homebrew/homebrew).


## Installing Helm Classic

From a Linux or Mac OS X client:
```
curl -s https://get.helm.sh | bash
```

*or*:

1. Grab a prebuilt binary from:
  - the latest release: [ ![Download](https://api.bintray.com/packages/deis/helm/helmc/images/download.svg) ](https://bintray.com/deis/helm/helmc/_latestVersion#files)
  - the CI build pipeline: [ ![Download](https://api.bintray.com/packages/deis/helm-ci/helmc/images/download.svg) ](https://bintray.com/deis/helm-ci/helmc/_latestVersion#files)
2. Unzip the package and make sure `helmc` is available on the PATH.

### Migration Notes

If you are a user of the original Helm tool (versions prior to v0.7.0), take note that Helm Classic is a _re-branding_ of that tool-- new name, same great taste!

__Helm Classic is fully compatible with previously existing Helm charts!__

Anyone migrating to Helm Classic from an older version should complete the following steps to fully replace their existing tool with Helm Classic. Doing so will clear the path for eventual installation of the new and improved Helm ([kubernetes/helm][k8s-helm]).

First, you may optionally define a custom home directory for use by Helm Classic. If opting for this, the instruction should be added to your shell's profile.

```
$ HELMC_HOME=/custom/path
```

Next, we migrate the contents from its old location to its new location (because the default location has changed).

```
$ mv $(helm home) $(helmc home)
```

Finally, remove the old version:

```
$ rm $(which helm)
```

You may now use the new binary, `helmc`, just as you previously used `helm`.  Soon, the `helm` name will be taken over by the new and improved Helm ([kubernetes/helm][k8s-helm]) and you will be able to make use of `helmc` in parallel with `helm` for as long as you have that need.

### Prerequisite

Helm Classic requires an appropriately wired `kubectl` client to speak with a running Kubernetes cluster.

## Using Helm Classic

To quickly install a redis cluster:

```
$ helmc update
---> Cloning into '$HOME/.helmc/cache/charts'...
---> Updating cache from https://github.com/helmc/charts
---> Done
$ helmc search redis
---> 	redis-cluster (redis-cluster 0.0.5) - Highly available Redis cluster with multiple sentinels and standbys.
---> 	redis-standalone (redis-standalone 0.0.1) - Standalone Redis Master
$ helmc info redis-cluster
Name: redis-cluster
Home: http://github.com/deis/redis-cluster
Version: 0.0.5
Description: Highly available Redis cluster with multiple sentinels and standbys.
Details: This package provides a highly available Redis cluster with multiple sentinels and standbys. Note the `redis-master` pod is used for bootstrapping only and can be deleted once the cluster is up and running.
$ helmc install redis-cluster
---> No chart named "redis-cluster" in your workspace. Fetching now.
---> Fetched chart into workspace $HOME/.helmc/workspace/charts/redis-cluster
---> Running `kubectl create -f` ...
services/redis-sentinel
pods/redis-master
replicationcontrollers/redis
replicationcontrollers/redis-sentinel
---> Done
```

To fetch, modify and install a chart out of your local workspace:

```
$ helmc update
---> Updating cache from https://github.com/helm/charts
---> Done
$ helmc fetch redis-standalone redis
---> Fetched chart into workspace $HOME/.helmc/workspace/charts/redis
---> Done
$ helmc edit redis
$ helmc install redis
---> Running `kubectl create -f` ...
replicationcontrollers/redis-standalone
---> Done
```

## Building the Helm Classic CLI

- Make sure you have a `kubectl` client installed and configured to speak with a running Kubernetes cluster.
- Helm Classic utilizes a containerized build and test process. Making use of the containerized development environment only requires `make` and a _local_ Docker daemon _or_ `docker-machine`.

```console
git clone https://github.com/helm/helm-classic.git

cd helm-classic

make bootstrap # installs all of helm classic's dependencies

make build # generates bin/helmc binary

./bin/helmc # prints usage

# optional

make install # installs helmc system-wide

helmc # prints usage

```

## License

Copyright 2015 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

[k8s-helm]: https://github.com/kubernetes/helm
