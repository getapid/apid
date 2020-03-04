# expect

## Summary

Expect are the validations to be done on the response.

## Fields

### code

{{ field\(type="int", required="no", desc="The expected status code

### body

#### body.type

This affect how some other fields are interpreted, such as `body.exact`. Currently only JSON and plaintext are supported.

#### body.content

A string with the expected response body. See `body.exact` about more details

#### body.exact

Whether or not to strictly validate the response body. If `false` and `body.type=json`, just the fields are recursively validated, but not scalar \(ints, strings, etc.\); for arrays, the fields of each element of the response array are validated recursively against the first element in the `expect` array. If `false` and `body.type=plaintext`, the the response needs to contain the `body.content`, but doesn't have to fully match it.

## Examples

### Exact JSON

```yaml
expect:
  code: 200
  body:
    type: 'json'
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

  
 In this case an API response below on the left will pass validation, but the one of the right will not \("Boris" != "Bobby"\)   
  


 \`\`\`json { "first\_name": "Bobby", "last\_name": "Hounslow", "address": { "postcode": "TW4 7AE" } } \`\`\`

 \`\`\`json { "first\_name": "Boris", "last\_name": "Hounslow", "address": { "postcode": "TW4 7AE" } } \`\`\`

&lt;/div&gt;

### Non-exact JSON

```yaml
expect:
  code: 200
  body:
    type: 'json'
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

  
 In this case an API response below on the left will pass validation, but the one of the right won't \("code" != "postcode"\)   
  


 \`\`\`json { "first\_name": "John", "last\_name": "Leicester", "address": { "postcode": "LE9 6HF" } } \`\`\` \`\`\`json { "first\_name": "John", "last\_name": "Leicester", "address": { "code": "LE9 6HF" } } \`\`\`

