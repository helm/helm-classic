# Using Other Repositories

Helm allows for the use of additional (potentially private) repositories of charts via the `helm repo` command.

## Adding a repository

`$ helm repo add mycharts https://github.com/dev/mycharts` will add a chart table with the name `mycharts` pointing to the `dev/mycharts` git repository (any valid git protocol with regular git authentication).

## Listing repositories

```
$ helm repo list
    charts*    https://github.com/helm/charts
    mycharts    https://github.com/dev/mycharts
```
Note the `*` indicates the default repository. This is configured in a `config.yaml` file in `$HELM_HOME`.

## Using a different repository

`$ helm fetch mycharts/app` will fetch the `app` chart from the `mycharts` repo. I can then `helm install` as normal.

## Removing repositories

`$ helm repo rm mycharts` Note: there is no confirmation requested.
