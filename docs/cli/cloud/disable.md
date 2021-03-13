# Disable

Disable removes the suite from the execution index so it won't be executed again until enabled.

## Flags

| Flag   | Short | Required | Default | Description                                                                                                                   |
| :----- | :---- | :------- | :------ | :---------------------------------------------------------------------------------------------------------------------------- |
| --key  | -k    | yes      |         | [Your API key](../cloud/README.md), this can also be injected via the `APID_KEY` environment variable. Flag takes precedence. |
| --name | -n    | yes      |         | The name of the suite                                                                                                         |

# Examples

```bash
apid cloud suite disable -k <apid cloud key> -n simple-suite
```
