# Transactions and steps

An APId configuration file consists of [transactions](../transactions) which in turn consist of [steps](../steps). Steps
are the basic elements of the configuration. They specify how to make a request and then how to validate
its response. Transactions bundle steps together to help you represent meaningful stories.

# Variables

APId allows you to have [variables](../variables) that will be inplaced throughout your steps. Any string value in
the yaml config can contain templates. Variables can be declared for the transaction or step scope or be global.
They can also come from the environment, which can be handy for things like secrets and passwords, or from a
response from your API.

# Examples

```yaml
version: 1
variables:
  api_url: 'https://jsonplaceholder.typicode.com'
transactions:
  - id: 'transaction-one'
    variables:
      title: 'delectus aut autem'
      body: 'quia et suscipit suscipit recusandae consequuntur expedita'
    steps:
      - id: 'authenticate'
        request:
          method: POST
          endpoint: '{{ var.api_url }}/auth'
          body: '{{ env.USERNAME }}:{{ env.PASSWORD }}'
        export:
          auth_header: 'response.headers.X-APIDAUTH'

      - id: 'get first todo'
        request:
          method: GET
          endpoint: '{{ var.api_url }}/todos/1'
          headers:
            X-APIDAUTH: '{{ authenticate.auth_header }}'
        expect:
          code: 200
          body:
            type: 'json'
            exact: true
            content: |
              {
                "title": "{{ var.title }},
                "id": 1,
                "text": "{{ var.body }}",
                "completed": false
              }
```
