local apid = import 'apid/apid.libsonnet';
local is = import 'apid/is.libsonnet';

local vars = import 'vars.libsonnet';

{
  ['%s-%s-%s' % [method, body, expectedBody]]: apid.spec(
    steps=[
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
    ],
  )
  for method in ['POST', 'PUT', 'PATCH', 'DELETE']
  for body in ['this is the body']
  for expectedBody in [
    is.regex('.+'),
    "this is the body",
    is.regex('this is \\w+ body'),
  ]
}
