local vars = import 'vars.libsonnet';

local text_body_spec(method, body, expectedBody) = std.manifestJson(
  {
    steps: [
      {
        name: "first request",
        request: {
          method: method,
          url: vars.url,
          body: body,
        },
        expect: {
          code: 200,
          body: expectedBody,
        },
      },
      {
        name: "second request",
        request: {
          method: method,
          url: vars.url,
          body: body,
        },
        expect: {
          code: 200,
          body: expectedBody,
        },
      }
    ]
  }
);

{
  ["%s-%s-%s" % [method, "body", expectedBody]]: text_body_spec(method, body, expectedBody) 
    for method in ["POST", "PUT", "PATCH", "DELETE"]
    for body in ["this is the body"]
    for expectedBody in ["", "is the b", "afdglajga"] 
}