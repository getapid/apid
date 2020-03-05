# check

## Check

Check will run all the transactions defined in `apid.yaml` in the current directory or optionally take a path to the config via the `-c` or `--config` flag. To get familiar with the syntax of a config file, see [Reference](https://github.com/getapid/apid-cli/tree/c25493e27ca1c4680ab6be23887bbb3e71fff850/docs/reference/README.md)

## Details

Check takes all the transactions that you have specified in the yaml. The steps in each transaction are executed sequentially. If a step fails, then the whole transaction is aborted and the rest of the steps are ignored. If a transaction fails, this will be reported in the console and the next transaction will be started.

## Examples

```bash
apid check
apid check --config ./tests/e2e/apid.yaml
```

