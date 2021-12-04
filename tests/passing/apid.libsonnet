local type_matcher(type) =
  {
    '$$matcher_type$$': 'type::%s' % type,
    '$$matcher_params$$': null,
  };

local spec(steps) = std.manifestJson(
  {
    steps: steps,
  },
);

local type = {
  int: type_matcher('int'),
  float: type_matcher('float'),
  string: type_matcher('string'),
  object: type_matcher('object'),
  array: type_matcher('array'),
};

local key(matcher) = '%s%s' % ['$$shorthand_matcher$$', std.manifestJsonEx(matcher, '')];
local any() = {
  '$$matcher_type$$': 'any',
};

local string(string, case_sensitive=true) = {
  '$$matcher_type$$': 'string',
  '$$matcher_params$$': {
    value: string,
    case_sensitive: case_sensitive,
  },
};

local regex(regex) = {
  '$$matcher_type$$': 'regex',
  '$$matcher_params$$': regex,
};

local int(int) = {
  '$$matcher_type$$': 'int',
  '$$matcher_params$$': int,
};

local float(float) = {
  '$$matcher_type$$': 'float',
  '$$matcher_params$$': float,
};

local json(json, subset=false) = {
  '$$matcher_type$$': 'json',
  '$$matcher_params$$': {
    value: json,
    subset: subset,
  },
};
local array(array, subset=false) = {
  '$$matcher_type$$': 'array',
  '$$matcher_params$$': {
    value: array,
    subset: subset,
  },
};

local len(len) = {
  '$$matcher_type$$': 'len',
  '$$matcher_params$$': len,
};

local and(matchers) = {
  '$$matcher_type$$': 'and',
  '$$matcher_params$$': matchers,
};

local or(matchers) = {
  '$$matcher_type$$': 'or',
  '$$matcher_params$$': matchers,
};

local range(from, to) = {
  '$$matcher_type$$': 'range',
  '$$matcher_params$$': {
    from: from,
    to: to,
  },
};
{}
