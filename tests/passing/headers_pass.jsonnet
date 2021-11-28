local is = import 'apid/is.libsonnet';

local vars = import 'vars.libsonnet';

local header_spec(method, name, value, expectedName, expectedValue) = std.manifestJson(
  {
    steps: [
      {
        name: 'first request',
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
            [expectedName]: expectedValue,
          },
        },
      },
      {
        name: 'second request',
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
            [expectedName]: expectedValue,
          },
        },
      },
    ],
  }
);

{
  ['%s-%s-%s-%s-%s' % [method, name, value, expectedName, expectedValue]]: header_spec(method, name, value, expectedName, expectedValue)
  for method in ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']
  for name in ['COOKIES']
  for expectedName in [is.key.string('cookies', false), is.key.regex('C\\wokies'), is.key.regex('[cC]ookies'), is.key.regex('\\w+')]
  for value in ['COOKIE_VALUE']
  for expectedValue in ['COOKIE_VALUE', is.string('CookIe_Value', false), is.regex('\\w+_\\w+')]
} + {
  ['%s-%s-%s-%s-%s' % [method, name, value, expectedName, expectedValue]]: header_spec(method, name, value, expectedName, expectedValue)
  for method in ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']
  for name in ['Authorization']
  for expectedName in ['Authorization', is.key.string('authorization', false)]
  for value in ['Bearer 469db547-65a0-4745-a0ac-0f821ca915d7']
  for expectedValue in [is.string('bearer 469db547-65a0-4745-a0ac-0f821ca915d7', false), is.regex('Bearer \\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b'), is.regex('\\w+ \\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b')]
}
