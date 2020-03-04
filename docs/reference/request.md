# Summary

Request specifies how to make the request - which endpoint to go to, what body to use, etc.

# Fields

## method

{{ field(type="string", required="yes", desc="The HTTP method of the request

## endpoint

{{ field(type="string", required="yes", desc="The complete URL of the request

## headers

{{ field(type="mapping", required="no", desc="Headers to attach to the request. Keys may repeat. If keys repeat, all the values are added to the header.

## body

{{ field(type="string", required="no", desc="A string of the body of the request

## skip_ssl_verify

{{ field(type="bool", required="no", default="false" desc="Whether or not to ignore certificate errors

# Examples

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
