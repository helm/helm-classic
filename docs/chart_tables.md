# Using Other Repositories

Helm Classic allows for the use of additional (potentially private) repositories of charts via the `helmc repo` command.

## Adding a repository

`$ helmc repo add mycharts https://github.com/dev/mycharts` will add a chart table with the name `mycharts` pointing to the `dev/mycharts` git repository (any valid git protocol with regular git authentication).

## Listing repositories

```
$ helmc repo list
    charts*    https://github.com/helm/charts
    mycharts    https://github.com/dev/mycharts
```
Note the `*` indicates the default repository. This is configured in a `config.yaml` file in `$HELMC_HOME`.

## Using a different repository

`$ helmc fetch mycharts/app` will fetch the `app` chart from the `mycharts` repo. I can then `helmc install` as normal.

## Removing repositories

`$ helmc repo rm mycharts` Note: there is no confirmation requested.
