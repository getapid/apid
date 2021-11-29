package matcher

import (
	"fmt"
)

type typeFloatMatcher struct{}

func TypeFloatMatcherWithOptions(value interface{}) Matcher {
	return &typeFloatMatcher{}
}

func (m typeFloatMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch data.(type) {
	case int64:
		return true, []string{fmt.Sprintf("%s is float", location)}, nil
	case float64:
		return true, []string{fmt.Sprintf("%s is float", location)}, nil
	}

	return false, nil, []string{fmt.Sprintf("%s wanted float, got %v", location, data)}
}

func (m typeFloatMatcher) String() string {
	return "type::float"
}
