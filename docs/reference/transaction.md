# transaction

## Summary

A transaction is a list of [steps](step/README.md) which are executed sequentially. If a step fails, the whole transaction fails.

## Fields

### id

{{ field\(type="string", required="yes", desc="A string to uniquely identify a transaction

### variables

{{ field\(type="[`variables`](variables/README.md)", required="no", desc="Variables scoped to this transaction

### steps

{{ field\(type="[`[]step`](step/README.md)", required="yes", desc="A list of steps

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
