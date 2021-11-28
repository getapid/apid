local apid = import 'apid/apid.libsonnet';
local is = import 'apid/is.libsonnet';

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
      'random float': is.float(66.861),
      [is.key.string('random')]: is.int(88),
      [is.key.regex('first\\w+')]: 'Lilith',
      Stephanie: {
        age: 93,
        address: is.json({
          city: 'Kobe',
          country: 'Australia',
          countryCode: 'VE',
        }),
      },
      array: [
        'Marline',
        'Catharine',
      ],
    },
  ]
}
//  + {
//   ['string-%s-%s-%s' % [method, 'body', expected]]: json_body_spec(method, body, expected)
//   for method in ['POST', 'PUT', 'PATCH', 'DELETE']
//   for body in [vars.json]
//   for expected in [
//     [
//       {
//         selector: 'firstname',
//         is: 'Lilith',
//       },
//     ],
//   ]
// } + {
//   ['json-%s-%s-%s' % [method, 'body', expected]]: json_body_spec(method, body, expected)
//   for method in ['POST', 'PUT', 'PATCH', 'DELETE']
//   for body in [vars.json]
//   for expected in [
//     [
//       {
//         selector: 'Stephanie',
//         is: {
//           age: 93,
//         },
//       },
//     ],
//   ]
// } + {
//   ['array-%s-%s-%s' % [method, 'body', expected]]: json_body_spec(method, body, expected)
//   for method in ['POST', 'PUT', 'PATCH', 'DELETE']
//   for body in [vars.json]
//   for expected in [
//     [
//       {
//         selector: 'array of objects',
//         is: [
//           {
//             index: 0,
//             'index start at 5': 5,
//           },
//           {
//             index: 1,
//             'index start at 5': 6,
//           },
//           {
//             index: 2,
//             'index start at 5': 7,
//           },
//         ],
//       },
//     ],
//   ]
// } + {
//   ['unordered-array-%s-%s-%s' % [method, 'body', expected]]: json_body_spec(method, body, expected)
//   for method in ['POST', 'PUT', 'PATCH', 'DELETE']
//   for body in [vars.json]
//   for expected in [
//     [
//       {
//         selector: 'array of objects',
//         is: [
//           {
//             index: 2,
//             'index start at 5': 7,
//           },
//           {
//             index: 0,
//             'index start at 5': 5,
//           },
//           {
//             index: 1,
//             'index start at 5': 6,
//           },
//         ],
//       },
//     ],
//   ]
// } + {
//   ['array-index-array-%s-%s-%s' % [method, 'body', expected]]: json_body_spec(method, body, expected)
//   for method in ['POST', 'PUT', 'PATCH', 'DELETE']
//   for body in [vars.json]
//   for expected in [
//     [
//       {
//         selector: 'array of objects.1',
//         is: {
//           index: 1,
//           'index start at 5': 6,
//         },
//       },
//     ],
//   ]
// } + {
//   ['regex-simple-%s-%s-%s' % [method, 'body', expected]]: json_body_spec(method, body, expected)
//   for method in ['POST', 'PUT', 'PATCH', 'DELETE']
//   for body in [vars.json]
//   for expected in [
//     [
//       {
//         selector: 'regEx',
//         is: 'hello+ to you',
//       },
//     ],
//   ]
// } + {
//   ['regex-email-%s-%s-%s' % [method, 'body', expected]]: json_body_spec(method, body, expected)
//   for method in ['POST', 'PUT', 'PATCH', 'DELETE']
//   for body in [vars.json]
//   for expected in [
//     [
//       {
//         selector: 'email uses current data',
//         is: '^\\S+@\\S+\\.\\S+$',
//       },
//     ],
//   ]
// }



