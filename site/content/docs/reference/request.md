+++
title = "request"
description = "the block specifying what request to make"
template = "docs/article.html"
sort_by = weight
weight = 180
+++


{{ h2(text="Summary") }}

Request specifies how to make the request - which endpoint to go to, what body to use, etc.

{{ h2(text="Fields") }}

{{ h3(text="method") }}

{{ field(type="string", required="true", desc="The HTTP method of the request") }}

{{ h3(text="endpoint") }}

{{ field(type="string", required="true", desc="The complete URL of the request") }}

{{ h3(text="headers") }}

{{ field(type="[headers](../headers)", required="false", desc="Headers to attach to the request") }}

{{ h3(text="body") }}

{{ field(type="string", required="false", desc="A string of the body of the request") }}

{{ h3(text="skip_ssl_verify") }}

{{ field(type="bool", required="false", default="false" desc="Whether or not to ignore certificate errors") }}

{{ h2(text="Examples") }}

```yaml
request:
  method: "GET"
  endpoint: "http://https://jsonplaceholder.typicode.com/posts/1"
  headers:
    X-APIDAUTH: "{{ env.AUTH_TOKEN }}"
```