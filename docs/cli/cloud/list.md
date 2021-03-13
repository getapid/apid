# List

List displays all the suites associated with the provided API key.

## Flags

| Flag   | Short | Required | Default | Description                                                                                                                   |
| :----- | :---- | :------- | :------ | :---------------------------------------------------------------------------------------------------------------------------- |
| --key  | -k    | yes      |         | [Your API key](../cloud/README.md), this can also be injected via the `APID_KEY` environment variable. Flag takes precedence. |
| --name | -n    | yes      |         | The name of the suite                                                                                                         |

# Examples

```bash
apid cloud suite list -k <apid cloud key>
```
