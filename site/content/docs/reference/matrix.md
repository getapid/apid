+++
title = "transaction.matrix"
description = "running the same transaction with different variables"
template = "docs/article.html"
sort_by = weight
weight = 205
+++

{{ h2(text="Summary") }}

A matrix allows running the same transaction with different variables. Different combinations of the variables will be
generated and the transaction will be ran with all of them. The order in which they are generated is not guaranteed.

{{ h2(text="Examples") }}

The below will send four different requests in four different transactions.

```yaml
id: transaction
matrix:
  api_url:
    - "http://localhost:8080"
    - "https://jsonplaceholder.typicode.com"
  todo_id:
    - 1
    - 2
steps:
  - id: todos-1
    request:
      method: GET
      endpoint: "{{ var.api_url }}/todos/{{ var.todo_id }}"
```

the different sets:
* `transaction-1: GET http://localhost:8080/todos/1`
* `transaction-2: GET http://localhost:8080/todos/2`
* `transaction-3: GET https://jsonplaceholder.typicode.com/todos/1`
* `transaction-4: GET https://jsonplaceholder.typicode.com/todos/2`
