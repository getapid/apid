+++
title = "check"
description = "check command"
template = "docs/article.html"
weight = 2
sort_by = "weight"
+++

{{ h2(text="Summary") }}

Check will run all the transactions defined in `apid.yaml` in the current directory or optionally take a path 
to the config via the `-c|--config` flag. To get familiar with the syntax of a config file, see [Reference](../../reference)

{{ h3(text="Details") }}

Check takes all the transactions that you have specified in the yaml. The steps in each transaction are executed
sequentially. If a step fails, then the whole transaction is aborted and the rest of the steps are ignored.
If a transaction fails, this will be reported in the console and the next transaction will be started.

{{ h3(text="Examples") }}

```shell script
apid check
apid check --config ./tests/e2e/apid.yaml
```
