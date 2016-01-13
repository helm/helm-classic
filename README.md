# Helm - The Kubernetes Package Manager

[![Build Status](https://travis-ci.org/helm/helm.svg?branch=master)](https://travis-ci.org/helm/helm) [![Go Report Card](http://goreportcard.com/badge/helm/helm)](http://goreportcard.com/report/helm/helm)

[Helm](https://helm.sh) bootstraps your Kubernetes cluster with **Charts** that provide ready-to-use workloads like:

- A Redis cluster
- A Postgres database
- An HAProxy edge load balancer

A Chart is a unit of Kubernetes manifests that reflect best practices as determined by the Helm community.  Helm's [architecture](docs/architecture.md) is heavily influenced by [Homebrew](https://github.com/Homebrew/homebrew).

To view or contribute Charts, head over to the [charts repo](https://github.com/helm/charts).

## Work in Progress

![Deis Graphic](https://s3-us-west-2.amazonaws.com/get-deis/deis-graphic-small.png)

Helm is changing quickly. Your feedback and participation are more than welcome, but be aware that this project is considered a work in progress.

Please note that Helm is not a general-purpose tool for managing a Kubernetes cluster.  For that, we recommend using [kubectl](http://kubernetes.io/v1.0/docs/user-guide/kubectl/kubectl.html).

### Updating from Helm 0.1

**If you are a Helm 0.1 user** you will need to do an extra step when
you upgrade to Helm 0.2 or later. We changed our GitHub org from `deis` to
`helm`, which means the new default charts repository is now
`github.com/helm/charts`.

To fix quickly, simply run a couple helm commands:

```
$ helm repo rm
$ helm repo add charts https://github.com/helm/charts
```

## Installing Helm

From a Linux or Mac OS X client:
```
curl -s https://get.helm.sh | bash
```

*or*:

1. Grab a prebuilt binary from:
  - the latest release: [ ![Download](https://api.bintray.com/packages/deis/helm/helm/images/download.svg) ](https://bintray.com/deis/helm/helm/_latestVersion#files)
  - the CI build pipeline: [ ![Download](https://api.bintray.com/packages/deis/helm-ci/helm/images/download.svg) ](https://bintray.com/deis/helm-ci/helm/_latestVersion#files)
2. Unzip the package and make sure `helm` is available on the PATH.

### Prerequisite

Helm requires an appropriately wired `kubectl` client to speak with a running Kubernetes cluster.

## Using Helm

To quickly install a redis cluster:

```
$ helm update
---> Cloning into '$HOME/.helm/cache/charts'...
---> Updating cache from https://github.com/helm/charts
---> Done
$ helm search redis
---> 	redis-cluster (redis-cluster 0.0.5) - Highly available Redis cluster with multiple sentinels and standbys.
---> 	redis-standalone (redis-standalone 0.0.1) - Standalone Redis Master
$ helm info redis-cluster
Name: redis-cluster
Home: http://github.com/deis/redis-cluster
Version: 0.0.5
Description: Highly available Redis cluster with multiple sentinels and standbys.
Details: This package provides a highly available Redis cluster with multiple sentinels and standbys. Note the `redis-master` pod is used for bootstrapping only and can be deleted once the cluster is up and running.
$ helm install redis-cluster
---> No chart named "redis-cluster" in your workspace. Fetching now.
---> Fetched chart into workspace $HOME/.helm/workspace/charts/redis-cluster
---> Running `kubectl create -f` ...
services/redis-sentinel
pods/redis-master
replicationcontrollers/redis
replicationcontrollers/redis-sentinel
---> Done
```

To fetch, modify and install a chart out of your local workspace:

```
$ helm update
---> Updating cache from https://github.com/helm/charts
---> Done
$ helm fetch redis-standalone redis
---> Fetched chart into workspace $HOME/.helm/workspace/charts/redis
---> Done
$ helm edit redis
$ helm install redis
---> Running `kubectl create -f` ...
replicationcontrollers/redis-standalone
---> Done
```

## Future Plans

Helm provides package manager semantics for Kubernetes workloads.  This is a novel concept.  As a result, there are a number of features that are not yet implemented and others that are not yet designed.  Your feedback is always appreciated.

Some of the features we plan to tackle in the near future:

- [x] Working with External (or Private) Repositories [#118](https://github.com/helm/helm/issues/118)
- [ ] Mandatory Group Labels [#80](https://github.com/helm/helm/issues/80)
- [x] Linting for Charts [#96](https://github.com/helm/helm/issues/96)
- [x] End-to-End Testing of Charts [#4](https://github.com/helm/helm/issues/4)
- [x] Dry-Run Installs [#78](https://github.com/helm/helm/issues/78)
- [ ] Improved Dependency Resolution (based on service provider)
- [ ] Upgrading Charts

For more detailed information on the project roadmap, see the [GitHub milestones](https://github.com/helm/helm/milestones).

## Contributing to the Helm CLI

- Make sure you have a `kubectl` client installed and configured to speak with a running Kubernetes cluster.
- Helm requires Go 1.5
- Install [glide](https://github.com/Masterminds/glide) >= 0.8.2
- Run the following commands:

```console
git clone https://github.com/helm/helm.git $GOPATH/src/github.com/helm/helm

cd $GOPATH/src/github.com/helm/helm

make bootstrap # installs all of helm's dependencies

make build # generates bin/helm binary

./bin/helm # prints usage

# optional

make install # installs helm system-wide

helm # prints usage

```

## License

Copyright 2015 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
