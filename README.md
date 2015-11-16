# Helm - The Kubernetes Package Manager

[![Build Status](https://travis-ci.org/deis/helm.svg?branch=master)](https://travis-ci.org/deis/helm) [![Go Report Card](http://goreportcard.com/badge/deis/helm)](http://goreportcard.com/report/deis/helm)

[Helm](https://helm.sh) bootstraps your Kubernetes cluster with **Charts** that provide ready-to-use workloads like:

- A Redis cluster
- A Postgres database
- An HAProxy edge load balancer

A Chart is a unit of Kubernetes manifests that reflect best practices as determined by the Helm community.  Helm's [architecture](docs/architecture.md) is heavily influenced by [Homebrew](https://github.com/Homebrew/homebrew).

To view or contribute Charts, head over to the [charts repo](https://github.com/deis/charts).

## Work in Progress

![Deis Graphic](https://s3-us-west-2.amazonaws.com/get-deis/deis-graphic-small.png)

Helm is changing quickly. Your feedback and participation are more than welcome, but be aware that this project is considered a work in progress.

Please note that Helm is not a general-purpose tool for managing a Kubernetes cluster.  For that, we recommend using [kubectl](http://kubernetes.io/v1.0/docs/user-guide/kubectl/kubectl.html).

## Installing Helm

1. Grab a prebuilt binary from:
  - the latest release: [ ![Download](https://api.bintray.com/packages/deis/helm/helm/images/download.svg) ](https://bintray.com/deis/helm/helm/_latestVersion#files)
  - the CI build pipeline: [ ![Download](https://api.bintray.com/packages/deis/helm-ci/helm/images/download.svg) ](https://bintray.com/deis/helm-ci/helm/_latestVersion#files)
2. Unzip the package and make sure `helm` is available on the PATH.
3. Install a `kubectl` client and configure it to speak with a running Kubernetes cluster.

## Using Helm

To quickly install a standalone etcd instance:

```
$ helm update
---> Cloning into '$HOME/.helm/cache/charts'...
---> Updating cache from https://github.com/deis/charts
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
---> Updating cache from https://github.com/deis/charts
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

- [ ] Working with External (or Private) Repositories [#118](https://github.com/deis/helm/issues/118)
- [ ] Mandatory Group Labels [#80](https://github.com/deis/helm/issues/80)
- [ ] Linting for Charts [#96](https://github.com/deis/helm/issues/96)
- [ ] End-to-End Testing of Charts [#4](https://github.com/deis/helm/issues/4)
- [x] Dry-Run Installs [#78](https://github.com/deis/helm/issues/78)
- [ ] Improved Dependency Resolution (based on service provider)
- [ ] Upgrading Charts

For more detailed information on the project roadmap, see the [GitHub milestones](https://github.com/deis/helm/milestones).

## Contributing to the Helm CLI

- Make sure you have a `kubectl` client installed and configured to speak with a running Kubernetes cluster.
- Helm requires Go 1.5
- Install [glide](https://github.com/Masterminds/glide) >= 0.7.2
- Run the following commands:

```console
git clone https://github.com/deis/helm.git $GOPATH/src/github.com/deis/helm

cd $GOPATH/src/github.com/deis/helm

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
