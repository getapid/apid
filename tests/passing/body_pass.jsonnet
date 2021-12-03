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
  ['float-%s-%s-%s' % [method, 'body', expected]]: spec(
    steps=steps(method, body, expected),
  )
  for method in ['POST', 'PUT', 'PATCH', 'DELETE']
  for body in [vars.json]
  for expected in [
    {
      'random float': and([
        float(66.861),
        type.float,
      ]),
      [key(string('random'))]: int(88),
      [key(regex('first\\w+'))]: 'Lilith',
      [key(
        or([
          string('Stephanie'),
          len(9),
        ])
      )]: {
        age: and([
          range(90, 94),
          type.int,
        ]),
        address: and([
          json({
            city: 'Kobe',
            country: 'Australia',
            countryCode: 'VE',
          }),
          type.object,
        ]),
      },
      array: and(
        [
          type.array,
          [
            and([
              'Marline',
              type.string,
            ]),
            'Catharine',
          ],
        ]
      ),
      countryCode: and(
        [
          len(2),
          string('VE'),
        ]
      ),
    },
  ]
}
