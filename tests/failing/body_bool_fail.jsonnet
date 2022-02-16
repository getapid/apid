local vars = import 'vars.libsonnet';

local text_body_spec(method, body, expectedBody) = std.manifestJson(
  {
    steps: [
      {
        name: 'first request',
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
        name: 'second request',
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
    ],
  }
);

{
  ['true-bool-%s-%s-%s' % [method, body, expectedBody]]: text_body_spec(method, body, expectedBody)
  for method in ['POST', 'PUT', 'PATCH', 'DELETE']
  for body in [true]
  for expectedBody in [
    false,
  ]
}
{
  ['false-bool-%s-%s-%s' % [method, body, expectedBody]]: text_body_spec(method, body, expectedBody)
  for method in ['POST', 'PUT', 'PATCH', 'DELETE']
  for body in [false]
  for expectedBody in [
    true,
  ]
}
{
  ['non-bool-%s-%s-%s' % [method, body, expectedBody]]: text_body_spec(method, body, expectedBody)
  for method in ['POST', 'PUT', 'PATCH', 'DELETE']
  for body in ['false', 12, 13.1247, {}, []]
  for expectedBody in [
    type.bool,
    false,
    true,
  ]
}
