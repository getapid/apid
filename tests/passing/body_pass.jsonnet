local _ = import 'apid/apid.libsonnet';

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
  ['float-%s-%s-%s' % [method, 'body', expected]]: _.spec(
    steps=steps(method, body, expected),
  )
  for method in ['POST', 'PUT', 'PATCH', 'DELETE']
  for body in [vars.json]
  for expected in [
    {
      'random float': _.and([
        _.float(66.861),
        _.type.float(),
      ]),
      [_.key(_.string('random'))]: _.int(88),
      [_.key(_.regex('first\\w+'))]: 'Lilith',
      [_.key(
        _.or([
          _.string('Stephanie'),
          _.len(9),
        ])
      )]: {
        age: _.and([
          _.range(90, 94),
          _.type.int()
        ]),
        address: _.and([
          _.json({
            city: 'Kobe',
            country: 'Australia',
            countryCode: 'VE',
          }),
          _.type.object()
        ]),
      },
      array: _.and(
        [
          _.type.array(),
          [
            _.and([
              'Marline',
              _.type.string(),
            ]),
            'Catharine',
          ],
        ]
      ),
      countryCode: _.and(
        [
          _.len(2),
          _.string('VE'),
        ]
      ),
    },
  ]
}
