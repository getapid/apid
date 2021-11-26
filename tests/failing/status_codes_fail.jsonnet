local vars = import 'vars.libsonnet';

local status_code_spec(method, code) = std.manifestJson(
  {
    steps: [
      {
        name: "first request",
        request: {
          type: method,
          url: vars.url,
          headers: {
            "X-ECHO-STATUSCODE": "%s" % (code + 1)
          },
        },
        expect: {
          code: code
        },
      },
      {
        name: "second request",
        request: {
          type: method,
          url: vars.url,
          headers: {
            "X-ECHO-STATUSCODE": "%s" %  (code + 1)
          },
        },
        expect: {
          code: code
        },
      }
    ]
  }
);

{
  ["%s-%d" % [method, code]]: status_code_spec(method, code) for method in ["GET", "POST", "PUT", "PATCH", "DELETE"] for code in [200, 300, 400, 500]
}