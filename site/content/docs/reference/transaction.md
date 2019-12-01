+++
title = "transaction"
description = ""
template = "docs/article.html"
sort_by = weight
weight = 200
+++

{{ h2(text="Summary") }}

A transaction is a list of [steps](../step) which are executed sequentially. If a step fails, the whole
transaction is aborted.

{{ h2(text="Fields") }}
 
{{ h3(text="id") }}

{{ field(type="string", required="true", desc="A string to uniquely identify a transaction") }}
 
{{ h3(text="variables") }}

{{ field(type="[`variables`](../variables)", required="false", desc="Variables scoped to this transaction") }}
 
{{ h3(text="steps") }}

{{ field(type="[`[]step`](../step)", required="true", desc="A list of steps") }}

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