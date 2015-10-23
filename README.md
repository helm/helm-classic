# Helm - The Kubernetes Package Manager

[![Build Status](https://travis-ci.org/deis/helm.svg?branch=master)](https://travis-ci.org/deis/helm)

Helm helps you bootstrap your Kubernetes cluster with **Charts** that provide ready-to-use workloads like:

- A Postgres database
- A scale-out etcd cluster
- An HAProxy Edge Load Balancer

A Chart is a unit of Kubernetes manifests that reflect best practices as determined by the Helm community.  Helm's [architecture](docs/architecture.md) is heavily influenced by [Homebrew](https://github.com/Homebrew/homebrew).

## Work in Progress

![Deis Graphic](https://s3-us-west-2.amazonaws.com/get-deis/deis-graphic-small.png)

`helm` is changing quickly. Your feedback and participation are more than welcome, but be aware that this project is considered a work in progress.

## Install Helm

To just use Helm, you can run the following:

```console
$ go get github.com/deis/helm/...
```
## Using Helm

To use Helm you'll need a `kubectl` client installed and configured to speak with a running Kubernetes cluster.

To quickly install a standalone etcd instance:

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
---> Cached files into $HELM_WORKDIR/charts/etcd-standalone/
---> Running kubectl create -f ...
---> Done
```

To fetch, customize and install the same chart:

```
$ helm update
---> Fetching updates...
---> Done
$ helm fetch etcd-standalone etcd
---> Downloading etcd-standalone-2.2.0-beta3
---> Cached files into $HELM_WORKDIR/charts/etcd/
---> Done
$ helm edit etcd
$ helm install etcd
---> Found etcd in $HELM_WORKDIR/charts/etcd/ ...
---> Running kubectl create -f ...
```

## Contributing

Helm charts are developed as part of a community effort.  To contribute:

1. Fork Helm
2. `helm create mychart && helm edit mychart`
3. `helm publish mychart`
4. `cd $HELM_CACHEDIR && git checkout -b mychart`
5. `git add charts/mychart && git commit && git push`
6. Submit a Pull Request

For helm Go developers, we suggest installing
[glide](https://github.com/Masterminds/glde) and then running `make
build` or `make install`.

## License

Copyright 2015 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
