# Enable

Enable schedules a suite for execution with the predefined schedule at the predefined locations

## Flags

| Flag   | Short | Required | Default | Description                                                                                                                   |
| :----- | :---- | :------- | :------ | :---------------------------------------------------------------------------------------------------------------------------- |
| --key  | -k    | yes      |         | [Your API key](../cloud/README.md), this can also be injected via the `APID_KEY` environment variable. Flag takes precedence. |
| --name | -n    | yes      |         | The name of the suite                                                                                                         |

# Examples

```bash
apid cloud suite enable -k <apid cloud key> -n simple-suite
```
