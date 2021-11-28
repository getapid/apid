local apid = import 'apid/apid.libsonnet';
local _ = import 'apid/is.libsonnet';

local vars = import 'vars.libsonnet';

local steps(method, body, expected) = [
  {
    name: 'first request',
    request: {
      method: method,
      url: vars.url,
      body: body,
    },
    expect: {
      code: 200,
      body: expected,
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
      body: expected,
    },
  },
];

{
  ['float-%s-%s-%s' % [method, 'body', expected]]: apid.spec(
    steps=steps(method, body, expected),
  )
  for method in ['POST', 'PUT', 'PATCH', 'DELETE']
  for body in [vars.json]
  for expected in [
    {
      'random float': _.float(66.861),
      [_.key.string('random')]: _.int(88),
      [_.key.regex('first\\w+')]: 'Lilith',
      [_.key.or([
        _.string("Stephanie"),
         _.len(9)
         ])]: {
        age: 93,
        address: _.json({
          city: 'Kobe',
          country: 'Australia',
          countryCode: 'VE',
        }),
      },
      array: [
        'Marline',
        'Catharine',
      ],
      countryCode: _.and(
        [
          _.len(2),
          _.string('VE')
        ]
      )
    },
  ]
}