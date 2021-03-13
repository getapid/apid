# Upsert

Upsert will create a new suite in APId cloudcontaining all the transactions defined in `apid.yaml` in the current directory. If a suite with the specified name exists it will replace it instead.

It can optionally take a path to the config via the `-c` or `--config` flag. To get familiar with the syntax of a config file, see [here](../../cloud.md)

## Details

Suites in APId cloud use the same syntax as the locally ran configs, thus it is really easy to upload your existing tests with one command. The only difference is that you'll need to provide a couple of extra top level values - `schedule` and a list of `locations`. These are needed so APId knows when and where to execute your suite. You can learn more about these in the config [Reference](../reference/README.md) section.

## Flags

| Flag     | Short | Required | Default     | Description                                                                                                                   |
| :------- | :---- | :------- | :---------- | :---------------------------------------------------------------------------------------------------------------------------- |
| --name   | -n    | yes      |             | The name of the suite, used to later reference it (enable, replace)                                                           |
| --key    | -k    | yes      |             | [Your API key](../cloud/README.md), this can also be injected via the `APID_KEY` environment variable. Flag takes precedence. |
| --config | -c    | no       | ./apid.yaml | The config file. If a folder is provided will recursively load all `*.yaml` files                                             |

# Examples

```bash
apid cloud suite upsert -c ./tests/e2e/apid.yaml -n simple-suite -k <apid cloud key>
```
