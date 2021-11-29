package matcher

import (
	"fmt"
)

type typeIntMatcher struct{}

func TypeIntMatcherWithOptions(value interface{}) Matcher {
	return &typeIntMatcher{}
}

func (m typeIntMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch val := data.(type) {
	case int64:
		return m.check(float64(val), location)
	case float64:
		return m.check(val, location)
	}

	return false, nil, []string{fmt.Sprintf("%s wanted int, got %v", location, data)}
}

func (m typeIntMatcher) check(val float64, location string) (bool, []string, []string) {
	if val == float64(int(val)) {
		return true, []string{fmt.Sprintf("%s is int", location)}, nil
	}
	return false, nil, []string{fmt.Sprintf("%s wanted int, got %v", location, val)}
}

func (m typeIntMatcher) String() string {
	return "type::integer"
}
