+++
title = "check"
description = "Running the transactions in apid.yaml"
template = "docs/article.html"
weight = 2
sort_by = "weight"
+++

{{ h2(text="Summary") }}

Check will run all the transactions defined in `apid.yaml` in the current directory. It can optionally take a path 
via the `-c|--config` flag. If the path is a directory all yaml files in the directory will be loaded.
To get familiar with the syntax of a config file, see [Reference](../../reference)

{{ h3(text="Details") }}

`check` takes all the transactions that you have specified in the yaml. The steps in each transaction are executed
sequentially. If a step fails, then the whole transaction is aborted and the rest of the steps are ignored.
If a transaction fails, this will be reported in the console and the next transaction will be started.

{{ h3(text="Examples") }}

```sh
apid check
apid check --config ./tests/e2e/apid.yaml
```

The output should look something like:
```text
simple-transaction:
    OK   health-endpoint-test-get
    OK   echo-endpoint-test-post
successful transactions: 1/1
failed transactions:     0/1
```