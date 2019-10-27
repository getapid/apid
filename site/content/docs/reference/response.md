+++
title = "response"
description = "the block specifying what request to make"
template = "docs/article.html"
sort_by = weight
weight = 181
+++


{{ h2(text="Summary") }}

Response is a set of validations that can be done on the response in a step.

{{ h2(text="Fields") }}

{{ h3(text="code") }}

{{ field(type="int", required="true", desc="The expected status code") }}

{{ h2(text="Examples") }}

```yaml
response:
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
```