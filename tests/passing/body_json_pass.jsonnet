local vars = import 'vars.libsonnet';

local json_body_spec(method, body, expectedBody) = std.manifestJson(
  {
    steps: [
      {
        name: "first request",
        request: {
          method: method,
          url: vars.url,
          body: std.manifestJsonEx(body, ""),
        },
        expect: {
          code: 200,
          json: expectedBody,
        },
      },
      {
        name: "second request",
        request: {
          method: method,
          url: vars.url,
          body: std.manifestJsonEx(body, ""),
        },
        expect: {
          code: 200,
          json: expectedBody,
        },
      }
    ]
  }
);

{
  ["float-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "random float",
          is: 66.861
        }
      ]
    ] 
} + {
  ["int-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "random",
          is: 88
        }
      ]
    ] 
} + {
  ["string-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "firstname",
          is: "Lilith"
        }
      ]
    ] 
} + {
  ["json-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "Stephanie",
          is: {
            age: 93
          }
        }
      ]
    ] 
} + {
  ["array-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "array of objects",
          is: [
            {
              "index": 0,
              "index start at 5": 5
            },
            {
              "index": 1,
              "index start at 5": 6
            },
            {
              "index": 2,
              "index start at 5": 7
            }
          ]
        }
      ]
    ] 
} + {
  ["unordered-array-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "array of objects",
          is: [
            {
              "index": 2,
              "index start at 5": 7
            },
            {
              "index": 0,
              "index start at 5": 5
            },
            {
              "index": 1,
              "index start at 5": 6
            }
          ]
        }
      ]
    ] 
} + {
  ["array-index-array-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "array of objects.1",
          is: {
            "index": 1,
            "index start at 5": 6
          }
        }
      ]
    ] 
} + {
  ["regex-simple-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "regEx",
          is: "hello+ to you"
        }
      ]
    ] 
} + {
  ["regex-email-%s-%s-%s" % [method, "body", expectedBody]]: json_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in [vars.json]
    for expectedBody in [
      [
        {
          selector: "email uses current data",
          is: "^\\S+@\\S+\\.\\S+$"
        }
      ]
    ] 
}