# transaction

## Summary

A transaction is a list of [steps](https://github.com/getapid/apid-cli/tree/22534ec0dafbcd65c14c4b649fbab9b5f7ae7398/docs/step/README.md) which are executed sequentially. If a step fails, the whole transaction fails.

## Fields

### id

{{ field\(type="string", required="yes", desc="A string to uniquely identify a transaction

### variables

{{ field\(type="[`variables`](https://github.com/getapid/apid-cli/tree/22534ec0dafbcd65c14c4b649fbab9b5f7ae7398/docs/variables/README.md)", required="no", desc="Variables scoped to this transaction

### steps

{{ field\(type="[`[]step`](https://github.com/getapid/apid-cli/tree/22534ec0dafbcd65c14c4b649fbab9b5f7ae7398/docs/step/README.md)", required="yes", desc="A list of steps

## Examples

```yaml
id: 'transaction-one'
variables:
  api_url: 'https://jsonplaceholder.typicode.com'
steps:
  - id: 'todos-1'
    request:
      method: 'GET'
      endpoint: '{{ var.api_url }}/todos/1'
```

