# Check

The cloud check command uses exactly the same configs as `check`, but executes requests from the APId cloud. This allows you to test your API from around the world. Remote will run all the transactions defined in `apid.yaml` in the current directory or optionally take a path to the config via the `-c` or `--config` flag. To get familiar with the syntax of a config file, see [Reference](../reference/README.md)

# Details

Check takes all the transactions that you have specified in the yaml. The steps in each transaction are executed sequentially. If a step fails, then the whole transaction is aborted and the rest of the steps are ignored. Each step is executed on a remote server, in a location specified via the command line flags. If a transaction fails, this will be reported in the console and the next transaction will be started.

## Flags

| Flag          | Short | Required | Default     | Description                                                                                                                   |
| :------------ | :---- | :------- | :---------- | :---------------------------------------------------------------------------------------------------------------------------- |
| --key         | -k    | yes      |             | [Your API key](../cloud/README.md), this can also be injected via the `APID_KEY` environment variable. Flag takes precedence. |
| --region      | -r    | no       | washington  | The location to run the tests from, a list or regions can be found [here](../cloud/README.md)                                 |
| --config      | -c    | no       | ./apid.yaml | The config file. If a folder is provided will recursively load all `*.yaml` files                                             |
| --timings     | -t    | no       | false       | Display the request timings, like DNS lookup, TCP connect, TLS handshake, etc                                                 |
| --parallelism | -p    | no       | 10          | Number of transactions to be run in parallel                                                                                  |

# Examples

```bash
apid cloud check --key <access-key>
apid cloud check --config ./tests/e2e/apid.yaml --key <access-key>
apid cloud check --config ./tests/e2e/apid.yaml --key <access-key> --region washington
```
