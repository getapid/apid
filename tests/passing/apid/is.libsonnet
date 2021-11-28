local SHORTHAND_MATCHER_PREFIX = '$$shorthand_matcher$$';


local any_matcher() =
  {
    '$$matcher_type$$': 'any',
  };

local string_matcher(string, case_sensitive=true) =
  {
    '$$matcher_type$$': 'string',
    '$$matcher_params$$': {
      value: string,
      case_sensitive: case_sensitive,
    },
  };

local regex_matcher(regex) =
  {
    '$$matcher_type$$': 'regex',
    '$$matcher_params$$': regex,
  };

local int_matcher(int) =
  {
    '$$matcher_type$$': 'int',
    '$$matcher_params$$': int,
  };

local float_matcher(float) =
  {
    '$$matcher_type$$': 'float',
    '$$matcher_params$$': float,
  };

local json_matcher(json, subset=false) =
  {
    '$$matcher_type$$': 'json',
    '$$matcher_params$$': {
      value: json,
      subset: subset,
    },
  };

local array_matcher(array, subset=false) =
  {
    '$$matcher_type$$': 'array',
    '$$matcher_params$$': {
      value: array,
      subset: subset,
    },
  };

local json_to_string(json) = std.manifestJsonEx(json, '');

{
  key: {
    regex(regex):: '%s%s' % [SHORTHAND_MATCHER_PREFIX, json_to_string(regex_matcher(regex))],
    string(string, case_sensitive=true):: '%s%s' % [SHORTHAND_MATCHER_PREFIX, json_to_string(string_matcher(string, case_sensitive))],
    int(int):: '%s%s' % [SHORTHAND_MATCHER_PREFIX, json_to_string(int_matcher(int))],
    float(float):: '%s%s' % [SHORTHAND_MATCHER_PREFIX, json_to_string(float_matcher(float))],
  },

  any():: any_matcher(),
  regex(regex):: regex_matcher(regex),
  string(string, case_sensitive=true):: string_matcher(string, case_sensitive),
  int(int):: int_matcher(int),
  float(float):: float_matcher(float),
  json(json, subset=false):: json_matcher(json, subset),
  array(array, subset=false):: array_matcher(array, subset),
}
