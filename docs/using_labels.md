# Helm Classic Labels

Helm Classic is designed to take full advantage of [Kubernetes labels](http://kubernetes.io/v1.0/docs/user-guide/labels.html).

## What are Labels?

From the Kubernetes documentation on the [motivation for labels](http://kubernetes.io/v1.0/docs/user-guide/labels.html#motivation):

> Labels enable users to map their own organizational structures onto system objects in a loosely coupled fashion, without requiring clients to store these mappings.
>
> Service deployments and batch processing pipelines are often multi-dimensional entities (e.g., multiple partitions or deployments, multiple release tracks, multiple tiers, multiple micro-services per tier). Management often requires cross-cutting operations, which breaks encapsulation of strictly hierarchical representations, especially rigid hierarchies determined by the infrastructure rather than by users.

To learn more about how labels work, check out [label selectors](http://kubernetes.io/v1.0/docs/user-guide/labels.html#label-selectors)
in the Kubernetes documentation.

## Helm Classic Labels

### Group Label

Helm Classic uses the `group` label as a convention for organizing Charts.  Services which share the same `group` are able to find each other and communicate automatically.  Examples include:

 * frontend
 * api
 * data

Groups are user-defined and not included in the Chart repository.

### Provider Label

Helm Classic uses the `provider` label as a convention specifying the type of Service provided by a Chart. Examples include:

* etcd
* postgres
* s3

A Chart may have dependencies on specific `provider`(s).  Chart authors are responsible for ensuring the `provider` label works consistently across Charts.

### Mode Label

Helm Classic uses the `mode` label as a convention for specifying the operating mode of the service.  Examples include:

* standalone
* clustered
* discovery

Charts may have dependencies on the operating `mode` of another Chart.

### Heritage Label

All Helm Classic Charts include the label `heritage: helm`. This provides a
convenient and standard way to query which components in a Kubernetes
cluster trace to Helm Classic.

## Using Labels

In Kubernetes, labels are typically edited by hand and stored with manifests in a version control system.  Helm Classic makes it easier to use labels effectively using the `helmc` CLI.

### Label Workflow (Simple)

If you want to place a package into a `group` while installing it, pass the group as an argument to `helmc install`.

```
helmc install nginx --group=frontend
helmc install python --group=frontend
```

### Label Workflow (Advanced)

Use the `helmc label` command to apply arbitrary labels to Charts in your workspace.

```
helmc fetch nginx
helmc fetch python
helmc label nginx group=frontend other=label
helmc label python group=frontend other=label
helmc install nginx
helmc install python
```

Of course, you can always use `helmc edit` or your own editor to customize labels and other manifest data in your local workspace.
