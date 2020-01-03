+++
title = "step.expect"
description = "The block specifying what response you are expecting"
template = "docs/article.html"
sort_by = weight
weight = 192
+++


{{ h2(text="Summary") }}

Expect are the validations to be done on the response.

{{ h2(text="Fields") }}

{{ h3(text="code") }}

{{ field(type="int", required="no", desc="The expected HTTP status code") }}

{{ h3(text="body") }}

{{ h4(text="body.type") }}

{{ field(type="string", required="no", default="plaintext", desc="This affects how some other fields are interpreted, such as `body.exact`.
See below for details. Currently only `json` and `plaintext` are supported.") }}

{{ h4(text="body.content") }}

{{ field(type="string", required="no", desc="A string with the expected response body. See `body.exact` for more details") }}

{{ h4(text="body.exact") }}

{{ field(type="bool", required="no", default="true", desc="Whether or not to strictly validate the response body. 
If `false` and `body.type=json`, just the fields are recursively validated, but not scalar (ints, strings, etc.);
for arrays, the fields of each element of the response array are validated recursively against the first element in the array in `expect`. 
If `false` and `body.type=plaintext`, the the response needs to contain the `body.content`, but doesn't have to fully match it.") }}

{{ h2(text="Examples") }}

{{ h3(text="Exact JSON") }}

```yaml
expect:
  code: 200
  body:
    type: "json"
    exact: true
    content: |
      {
        "first_name": "Bobby",
        "last_name": "Hounslow",
        "address": {
            "postcode": "TW4 7AE"
        }
      }
```
<br>
In this case an API response below on the left will pass validation, but the one of the right will not ("Boris" != "Bobby")
<br><br>

<div class="columns">
<div class="column is-6">

```json
    {
      "first_name": "Bobby",
      "last_name": "Hounslow",
      "address": {
        "postcode": "TW4 7AE"
      }
    }
```
</div>

<div class="column is-6">

```json
    {
      "first_name": "Boris",
      "last_name": "Hounslow",
      "address": {
        "postcode": "TW4 7AE"
      }
    }
```
</div>

</div>

{{ h3(text="Non-exact JSON") }}

```yaml
expect:
  code: 200
  body:
    type: "json"
    exact: false
    content: |
      {
        "first_name": "Bobby",
        "last_name": "Hounslow",
        "address": {
            "postcode": "TW4 7AE"
        }
      }
```

<br>
In this case an API response below on the left will pass validation, but the one of the right won't ("code" != "postcode")
<br><br>

<div class="columns">
<div class="column is-6">

```json
    {
      "first_name": "John",
      "last_name": "Leicester",
      "address": {
        "postcode": "LE9 6HF"
      }
    }
```

</div>
<div class="column is-6">

```json
    {
      "first_name": "John",
      "last_name": "Leicester",
      "address": {
        "code": "LE9 6HF"
      }
    }
```

</div>
</div>


{{ h3(text="Non-exact JSON Array") }}

```yaml
expect:
  code: 200
  body:
    type: "json"
    exact: false
    content: |
      {
        "people": [
          {
            "first_name": "Bobby",
            "last_name": "Hounslow"
          }
        ]
      }
```

<br>
In this case an API response below on the left will pass validation, but the one of the right won't
<br><br>

<div class="columns">
<div class="column is-6">

```json
    {
      "people": [
        {
           "first_name": "John",
           "last_name": "Leicester"
        },
        {
           "first_name": "John",
           "last_name": "Leicester"
        }
      ]
    }
```

</div>
<div class="column is-6">

```json
     {
       "people": [
         {
           "first_name": "John",
           "last_name": "Leicester"
         },
         {
            "first_name": "John"

         }
       ]
     }
```

</div>
</div>