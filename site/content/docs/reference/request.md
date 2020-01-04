+++
title = "step.request"
description = "The block specifying what request to make"
template = "docs/article.html"
sort_by = weight
weight = 198
+++


{{ h2(text="Summary") }}

Request specifies how to make the request - which endpoint to go to, what body to use, etc.

{{ h2(text="Fields") }}

{{ h3(text="method") }}

{{ field(type="string", required="yes", desc="The HTTP method of the request") }}

{{ h3(text="endpoint") }}

{{ field(type="string", required="yes", desc="The complete URL of the request") }}

{{ h3(text="headers") }}

{{ field(type="mapping", required="no", desc="Headers to attach to the request. Keys may repeat. If keys repeat or the values are arrays, all the values are added to the header.") }}

{{ h3(text="body") }}

{{ field(type="string", required="no", desc="A string of the body of the request") }}

{{ h3(text="skip_ssl_verify") }}

{{ field(type="bool", required="no", default="false" desc="Whether or not to ignore certificate errors") }}

{{ h2(text="Examples") }}

```yaml
request:
  method: "GET"
  endpoint: "http://https://jsonplaceholder.typicode.com/posts/1"
  headers:
    X-APIDAUTH: "{{ env.AUTH_TOKEN }}"
    Accept: "application/json"
    Accept: "application/ld+json"
    Accept:
      - "application/json"
      - "application/ld+json"
```