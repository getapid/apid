# Check

Check will run all the transactions defined in `apid.yaml` in the current directory or optionally take a path to the config via the `-c` or `--config` flag. To get familiar with the syntax of a config file, see [Reference](../reference/README.md)

# Details

Check takes all the transactions that you have specified in the yaml. The steps in each transaction are executed sequentially. If a step fails, then the whole transaction is aborted and the rest of the steps are ignored. If a transaction fails, this will be reported in the console and the next transaction will be started.

## Flags

| Flag          | Short | Required | Default     | Description                                                                       |
| :------------ | :---- | :------- | :---------- | :-------------------------------------------------------------------------------- |
| --config      | -c    | no       | ./apid.yaml | The config file. If a folder is provided will recursively load all `*.yaml` files |
| --timings     | -t    | no       | false       | Display the request timings, like DNS lookup, TCP connect, TLS handshake, etc     |
| --parallelism | -p    | no       | 10          | Number of transactions to be run in parallel                                      |

# Examples

```bash
apid check
apid check --config ./tests/e2e/apid.yaml
```
