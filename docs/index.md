# Helm Classic: The Original Kubernetes Package Manager

[Helm Classic](https://helm.sh) helps you find and use software built for Kubernetes. The Helm Classic CLI uses Charts which not only contain metadata about software packages but also Pod, ReplicationController and Service manifests for Kubernetes. With a few Helm Classic commands you can quickly and easily deploy software packages like:

- Postgres
- etcd
- HAProxy
- redis

All of the Helm Classic charts live at [github.com/helm/charts](https://github.com/helm/charts). If you want to make your own charts we have a guide for [authoring charts](authoring_charts.md) as well.

## Additional Information

Learn more about Helm Classic's [architecture](architecture.md).

Find out how Helm Classic uses [Kubernetes labels](using_labels.md).

If you are authoring a chart that will be used by other apps, check out how Helm Classic [models services](modeling_services.md).

## Thanks

Helm Classic was inspired by [Homebrew](https://github.com/Homebrew/homebrew).
