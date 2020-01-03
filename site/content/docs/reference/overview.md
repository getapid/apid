+++
title = "Overview"
description = "The structure of the YAML configuration file"
template = "docs/article.html"
sort_by = weight
weight = 10
+++

{{ h3(text="Transactions and steps") }}

An APId configuration file consists of [transactions](../transactions) which in turn consist of [steps](../steps). Steps 
are the basic elements of the configuration. They specify how to make a request and then how to validate
its response. Transactions bundle steps together to help you represent meaningful stories.

{{ h3(text="Variables") }}

APId has [variables](../variables) that can be inplaced throughout you steps and transactions.
Variables can be declared for the transaction or step scope or be global.
They can also come from the environment, which can be handy for things like secrets and passwords, or from a
response from your API. Most string values in the yaml config can contain templates (`{{ }}`).

{{ h3(text="Commands") }}

Commands allow to execute shell commands (using your default `$SHELL`) and use the output inside the request. They can be used anywhere where templates
can. They use the `{% %}` delimiters. See [commands](../commands) for more details.

{{ h3(text="Examples") }}

```yaml
version: 1
variables:
  api_url: "https://jsonplaceholder.typicode.com"
transactions:
  - id: "transaction-one"
    variables:
      title: "delectus aut autem"
      body: "quia et suscipit suscipit recusandae consequuntur expedita"
    steps:
      - id: "authenticate"
        request:
          method: POST
          endpoint: "{{ var.api_url }}/auth"
          headers:
            Authorization: "Basic {% echo -n $USERNAME:$PASSWORD | base64 %}"
        export:
          auth_header: "response.headers.X-APIDAUTH"

      - id: "get first todo"
        request:
          method: GET
          endpoint: "{{ var.api_url }}/todos/1"
          headers:
            X-APIDAUTH: "{{ authenticate.auth_header }}"
        expect:
          code: 200
          body:
            type: "json"
            exact: true
            content: |
              {
                "title": "{{ var.title }},
                "id": 1,
                "text": "{{ var.body }}",
                "completed": false
              }
```