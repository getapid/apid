+++
title = "transaction"
description = "A sequence of steps"
template = "docs/article.html"
sort_by = weight
weight = 200
+++

{{ h2(text="Summary") }}

A transaction is a list of [steps](../step) which are executed sequentially. If a step fails, the whole
transaction fails and further steps are not executed.

{{ h2(text="Fields") }}
 
{{ h3(text="id") }}

{{ field(type="string", required="yes", desc="A string to uniquely identify a transaction") }}
 
{{ h3(text="variables") }}

{{ field(type="[`variables`](../variables)", required="no", desc="Variables scoped to this transaction. See [variables](../variables)") }}
 
{{ h3(text="steps") }}

{{ field(type="[`[]step`](../step)", required="yes", desc="A list of steps. See [step](../step)") }}

{{ h2(text="Examples") }}

```yaml
id: "transaction-one"
variables:
  api_url: "https://jsonplaceholder.typicode.com"
steps:
  - id: "todos-1"
    request:
      method: "GET"
      endpoint: "{{ var.api_url }}/todos/1"
```
