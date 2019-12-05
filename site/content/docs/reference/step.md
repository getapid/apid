+++
title = "step"
description = "A step is a call to a single endpoint"
template = "docs/article.html"
sort_by = weight
weight = 190
+++

{{ h2(text="Summary") }}

A step is a call to a single endpoint with optional validation of the response.

{{ h2(text="Fields") }}

{{ h3(text="id") }}

{{ field(type="string", required="yes", desc="A string to uniquely identify a step within a transaction.") }}

{{ h3(text="request") }}

{{ field(type="[`request`](../request)", required="yes", desc="The request to send") }}

{{ h3(text="response") }}

{{ field(type="[`response`](../response)", required="no", desc="Validation on the response") }}

{{ h3(text="export") }}

{{ field(type="mapping", required="no", desc="The variables to export; a mapping from variable names to JSON paths
to the value in the response. You can export values from headers and JSON bodies. See the examples below.") }}

{{ h3(text="variables") }}

{{ field(type="[`variables`](../variables)", required="no", desc="Variables scoped to this step") }}

{{ h2(text="Examples") }}

```yaml
steps:
  - id: "get user with id 1"
    variables:
      api_url: "http://localhost:80"
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
      auth_header: "response.headers.X-APIDAUTH"
      auth_token: "response.body.access_token"
```
