# deployment-extractor

When a BOSH deployment fails, you cannot retrieve the manifest with a `bosh download manifest`. This program will extract the manifest from a `create deployment` task by tapping into its debug logs. It expects the BOSH debug log from `stdin` and will print the manifest to `stdout`. See example usage below.

## usage

`bosh task <insert task id> --debug | deployment-extractor`

