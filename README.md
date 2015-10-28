# Helm - The Kubernetes Package Manager

[![Build Status](https://travis-ci.org/deis/helm.svg?branch=master)](https://travis-ci.org/deis/helm) [![Go Report Card](http://goreportcard.com/badge/deis/helm)](http://goreportcard.com/report/deis/helm)

Helm helps you bootstrap your Kubernetes cluster with **Charts** that provide ready-to-use workloads like:

- A Postgres database
- A scale-out etcd cluster
- An HAProxy Edge Load Balancer

A Chart is a unit of Kubernetes manifests that reflect best practices as determined by the Helm community.  Helm's [architecture](docs/architecture.md) is heavily influenced by [Homebrew](https://github.com/Homebrew/homebrew).

To view or contribute Charts, head over to the [charts
repo](https://github.com/deis/charts).

## Work in Progress

![Deis Graphic](https://s3-us-west-2.amazonaws.com/get-deis/deis-graphic-small.png)

`helm` is changing quickly. Your feedback and participation are more than welcome, but be aware that this project is considered a work in progress.

## Installing Helm

1. Grab a prebuilt binary from:
  - the latest release: [ ![Download](https://api.bintray.com/packages/deis/helm/helm/images/download.svg) ](https://bintray.com/deis/helm/helm/_latestVersion)
  - the CI build pipeline: [ ![Download](https://api.bintray.com/packages/deis/helm-ci/helm/images/download.svg) ](https://bintray.com/deis/helm-ci/helm/_latestVersion)
2. Unzip the package and make sure `helm` is available on the PATH.
3. Install a `kubectl` client and configure it to speak with a running Kubernetes cluster.

## Using Helm

To quickly install a standalone etcd instance:

```
$ helm update
---> Fetching updates...
---> Done
$ helm search
--->
=================
---> Available Charts
---> =================

---> 	.git - UNKNOWN
---> 	README.md - UNKNOWN
---> 	alpine (alpine-pod 0.0.1) - Simple pod running Alpine Linux.
---> 	example-todo (example-todo 0.0.6) - Example Todo application backed by Redis
---> 	redis-cluster (redis-cluster 0.0.5) - Highly available Redis cluster with multiple sentinels and standbys.
---> 	redis-standalone (redis-standalone 0.0.1) - Standalone Redis Master
--->
$ helm info example-todo
---> /Users/nickleli/.helm/cache/charts/example-todo/Chart.yaml
---> Chart: example-todo
---> Description: Example Todo application backed by Redis
---> Details: This chart contains an Example Todo application backed by Redis. It supports standalone and clustered (sentinel) Redis backends.
---> Version: 0.0.6
---> Website: http://github.com/deis/example-todo
---> Dependencies: []
$ helm install example-todo
---> No installed chart named "example-todo". Installing now.
---> Fetching
---> Done
```


## Contributing to the Helm CLI

- Make sure you have a `kubectl` client installed and configured to speak with a running Kubernetes cluster.
- Install [glide](https://github.com/Masterminds/glide)
- Run the following commands:

```console
git clone https://github.com/deis/helm.git $GOPATH/github.com/deis/helm

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
