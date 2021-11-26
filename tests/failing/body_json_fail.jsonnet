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
          is: 66.862
        }
      ],
      [
        {
          selector: "random float",
          is: "66.862"
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
          is: 89
        }
      ],
      [
        {
          selector: "random",
          is: "89"
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
          is: "Liliadth"
        }
      ],
      [
        {
          selector: "firstname",
          is: ""
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
          is: {}
        }
      ],
      [
        {
          selector: "Stephanie",
          is: 93
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
            "index": 12,
            "index start at 5": 6,
            "some other key": "one"
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
          is: "hello"
        }
      ]
    ] 
} 