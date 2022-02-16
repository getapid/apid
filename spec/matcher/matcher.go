package matcher

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/getapid/apid/log"
)

const (
	MATCHER_TYPE_KEY         = "$$matcher_type$$"
	MATCHER_PARAMS_KEY       = "$$matcher_params$$"
	SHORTHAND_MATCHER_PREFIX = "$$shorthand_matcher$$"
)

var (
	ErrNotAMatcher          error = errors.New("not_a_matcher")
	ErrMalformedMatcherType error = errors.New("malformed_matcher_type")
)

type Matcher interface {
	Match(interface{}, string) (bool, []string, []string)
	String() string
}

func tryGetMatcherFromMap(m map[string]interface{}) Matcher {
	t, params, err := getMatcherTypeAndParams(m)
	if err != nil {
		switch err {
		case ErrNotAMatcher:
			return nil
		case ErrMalformedMatcherType:
			log.L.Fatalf("malformed matcher type")
		}
	}

	switch t {
	case "any":
		return AnyMatcher()
	case "regex":
		return RegexMatcherWithOptions(params)
	case "string":
		return StringMatcherWithOptions(params)
	case "int":
		return IntMatcherWithOptions(params)
	case "float":
		return FloatMatcherWithOptions(params)
	case "bool":
		return BoolMatcherWithOptions(params)
	case "json":
		return JSONMatcherWithOptions(params)
	case "array":
		return ArrayMatcherWithOptions(params)
	case "len":
		return LenMatcherWithOptions(params)
	case "and":
		return AndMatcherWithOptions(params)
	case "or":
		return OrMatcherWithOptions(params)
	case "range":
		return RangeMatcherWithOptions(params)
	case "type::int":
		return TypeIntMatcherWithOptions(params)
	case "type::bool":
		return TypeBoolMatcherWithOptions(params)
	case "type::float":
		return TypeFloatMatcherWithOptions(params)
	case "type::string":
		return TypeStringMatcherWithOptions(params)
	case "type::object":
		return TypeObjectMatcherWithOptions(params)
	case "type::array":
		return TypeArrayMatcherWithOptions(params)
	default:
		log.L.Fatalf("unknown matcher type %s", t)
		return nil
	}
}

func GetMatcher(i interface{}) Matcher {
	switch val := i.(type) {
	case map[string]interface{}:
		// this is a complex matcher
		if m := tryGetMatcherFromMap(val); m != nil {
			return m
		}
		return JSONMatcher(val)
	case []interface{}:
		return ArrayMatcher(val)
	case bool:
		return BoolMatcherWithOptions(val)
	case float64:
		return FloatMatcherWithOptions(val)
	case string:
		matcher := isShorthandMatcher(val)
		if matcher == nil {
			matcher = StringMatcher(val)
		}
		return matcher
	case nil:
		return NilMatcher()
	}
	// shouldn't get to here ever since we cover all cases
	// ref: https://pkg.go.dev/encoding/json#Unmarshal
	return nil
}

func getMatcherTypeAndParams(d map[string]interface{}) (string, interface{}, error) {
	t, ok := d[MATCHER_TYPE_KEY]
	if !ok {
		return "", nil, ErrNotAMatcher
	}

	var matcherType string
	switch ty := t.(type) {
	case string:
		matcherType = ty
	default:
		return "", nil, ErrMalformedMatcherType
	}

	return matcherType, d[MATCHER_PARAMS_KEY], nil
}

func isShorthandMatcher(str string) Matcher {
	if !strings.HasPrefix(str, SHORTHAND_MATCHER_PREFIX) {
		return nil
	}
	str = strings.TrimPrefix(str, SHORTHAND_MATCHER_PREFIX)
	var i map[string]interface{}
	err := json.Unmarshal([]byte(str), &i)
	if err != nil {
		log.L.Fatal("invalid shorthand matcher, %v", str)
	}
	m := tryGetMatcherFromMap(i)
	if m == nil {
		log.L.Fatal("invalid shorthand matcher, %v", str)
	}
	return m
}
