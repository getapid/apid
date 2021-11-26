local vars = import 'vars.libsonnet';

local header_spec(method, name, value, expectedName, expectedValue) = std.manifestJson(
  {
    steps: [
      {
        name: "first request",
        request: {
          method: method,
          url: vars.url,
          headers: {
            [name]: value,
          },
        },
        expect: {
          code: 200,
          headers: {
            [expectedName]: expectedValue
          },
        },
      },
      {
        name: "second request",
        request: {
          method: method,
          url: vars.url,
          headers: {
            [name]: value,
          },
        },
        expect: {
          code: 200,
          headers: {
            [expectedName]: expectedValue
          },
        },
      }
    ]
  }
);

{
  ["%s-%s-%s-%s-%s" % [method, name, value, expectedName, expectedValue]]: header_spec(method, name, value, expectedName, expectedValue) 
    for method in ["GET", "POST", "PUT", "PATCH", "DELETE"]
    for name in ["COOKIES"]
    for expectedName in ["cookies", "C\\wokies", "[cC]ookies", "\\w+"]
    for value in ["COOKIE_VALUE"]
    for expectedValue in ["COOKIE_VALUE", "CookIe_Value", "cookie_value", "\\w+_\\w+"] 
} + {
  ["%s-%s-%s-%s-%s" % [method, name, value, expectedName, expectedValue]]: header_spec(method, name, value, expectedName, expectedValue) 
    for method in ["GET", "POST", "PUT", "PATCH", "DELETE"]
    for name in ["Authorization"]
    for expectedName in ["Authorization", "authorization"]
    for value in ["Bearer 469db547-65a0-4745-a0ac-0f821ca915d7"]
    for expectedValue in ["bearer 469db547-65a0-4745-a0ac-0f821ca915d7", "Bearer \\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b", "\\w+ \\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b"]
}