package matcher

import (
	"fmt"
)

type typeArrayMatcher struct{}

func TypeArrayMatcherWithOptions(value interface{}) Matcher {
	return &typeArrayMatcher{}
}

func (m typeArrayMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch data.(type) {
	case []interface{}:
		return true, []string{fmt.Sprintf("%s is an array", location)}, nil
	}

	return false, nil, []string{fmt.Sprintf("%s wanted array, got %v", location, data)}
}

func (m typeArrayMatcher) String() string {
	return "type::array"
}
