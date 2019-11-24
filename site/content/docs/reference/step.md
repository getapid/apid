+++
title = "step"
description = ""
template = "docs/article.html"
sort_by = weight
weight = 190
+++

{{ h2(text="Summary") }}

A step is a call to a single endpoint. It may include some variables that are scoped only to it. You can also choose
to validate the response.

{{ h2(text="Fields") }}

{{ h3(text="id") }}

{{ field(type="string", required="true", desc="A string to uniquely identify a step within a transaction") }}

{{ h3(text="request") }}

{{ field(type="[`request`](../request)", required="true", desc="The request to send") }}

{{ h3(text="response") }}

{{ field(type="[`response`](../response)", required="false", desc="Validation on the response") }}

{{ h3(text="export") }}

{{ field(type="mapping", required="false", desc="The variables to export; a mapping from variable names to JSON paths
to the value in the response. See [variables](../variables) for more details") }}

{{ h3(text="variables") }}

{{ field(type="[`variables`](../variables)", required="false", desc="Variables scoped to this step") }}

{{ h2(text="Examples") }}

```yaml
steps:
  - id: "get user with id 1"
    request:
      method: "GET"
      endpoint: "{{ var.api_url }}/users/1"
    expect:
      code: 200
      body:
        type: "json"
        content: |
          {
            "first_name": "Bobby",
            "last_name": "Hounslow",
            "address": {
                "postcode": "TW4 7AE",
                "country": "UK",
            }
          }
    export:
      auth_token: "{{ response.headers.X-APIDAUTH }}"
```
