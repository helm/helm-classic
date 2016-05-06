# Helm Classic: A Kubernetes Package Manager

[Helm Classic](https://helm.sh) helps you find and use software built for Kubernetes. The Helm Classic CLI uses Charts which not only contain metadata about software packages but also Pod, ReplicationController and Service manifests for Kubernetes. With a few Helm Classic commands you can quickly and easily deploy software packages like:

- Postgres
- etcd
- HAProxy
- redis

All of the Helm Classic charts live at [github.com/helm/charts](https://github.com/helm/charts). If you want to make your own charts we have a guide for [authoring charts](authoring_charts.md) as well.

## Installing Helm Classic

From a Linux or Mac OS X client:
```
curl -s https://get.helm.sh | bash
```

*or*:

1. Grab a prebuilt binary from:
  - the latest release: [ ![Download](https://api.bintray.com/packages/deis/helm/helm-classic/images/download.svg) ](https://bintray.com/deis/helm/helmc/_latestVersion#files)
  - the CI build pipeline: [ ![Download](https://api.bintray.com/packages/deis/helm-ci/helm-classic/images/download.svg) ](https://bintray.com/deis/helm-ci/helmc/_latestVersion#files)
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

## Additional Information

Learn more about Helm Classic's [architecture](architecture.md).

Find out how Helm Classic uses [Kubernetes labels](using_labels.md).

If you are authoring a chart that will be used by other apps, check out how Helm Classic [models services](modeling_services.md).

## Thanks

Helm Classic was inspired by [Homebrew](https://github.com/Homebrew/homebrew).

[k8s-helm]: https://github.com/kubernetes/helm
