# Authoring Helm Charts

It is important for chart authors to understand Helm fundamentals.  Before you begin, make sure you are familiar with:

- How to [model Services in Helm](modeling_services.md)
- How Helm [uses Kubernetes Labels](using_labels.md)
- How the [Helm workspace](workspace.md) is laid out

## Background

Helm Charts consist of three items:

1. A `manifests` directory for Kubernetes resources
2. A `Chart.yaml` file
3. A `README.md`

The directory structure of a chart is as follows:

```
|- mychart/
    |
    |- manifests/
         |
         |- mychart-rc.yaml
         |- mychart-service.yaml
         |- ...
    |- Chart.yaml
    |- README.md
```

## Create a new Chart

### Step 1: Create the Chart in your Workspace

Use `helm create <chart-name>` to create a new chart in your workspace.
This will copy the default "skeleton" chart into `~/.helm/workspace/charts/<chart-name>`.

### Step 2: Edit the Chart

Use helm edit <chart-name> to open all files in the chart in a single editor.  

For convenience, this will present all the chart files inside a single editor, with `--- : <filepath>` delimiters.  This makes it easy to modify a chart, add files, and remove files all within a single `helm edit` command.

If you prefer to edit files manually, you can use an IDE or any other file-based editor.

### Step 3: Test the Chart

Use `helm test <chart-name>` to test installing the chart and validating that the proper Kubernetes resources are created, as evidenced by the `helm test` output and return code.

### Step 4: Publish the Chart

Use `helm publish <chart-name>` to copy a chart from your local workspace into the Git checkout that lives under `~/.helm/cache`.  From here you can submit a pull request.
