# Helm: The Kubernetes Package Manager

[Helm](https://helm.sh) helps you find and use software built for Kubernetes. The Helm CLI uses Charts which not only contain metadata about software packages but also Pod, ReplicationController and Service manifests for Kubernetes. With a few Helm commands you can quickly and easily deploy software packages like:

- Postgres
- etcd
- HAProxy
- redis

All of the Helm charts live at [github.com/deis/charts](https://github.com/deis/charts). If you want to make your own charts we have a guide for [authoring charts](authoring_charts.md) as well.

## Additional Information

Learn more about Helm's [architecture](architecture.md).

Find out how Helm uses [Kubernetes labels](using_labels.md).

If you are authoring a chart that will be used by other apps, check out how Helm [models services](modeling_services.md).

## Thanks

Helm was inspired by [Homebrew](https://github.com/Homebrew/homebrew).
