package matcher

import (
	"fmt"
)

type typeObjectMatcher struct{}

func TypeObjectMatcherWithOptions(value interface{}) Matcher {
	return &typeObjectMatcher{}
}

func (m typeObjectMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch data.(type) {
	case map[string]interface{}:
		return true, []string{fmt.Sprintf("%s is an object", location)}, nil
	}

	return false, nil, []string{fmt.Sprintf("%s wanted object, got %v", location, data)}
}

func (m typeObjectMatcher) String() string {
	return "type::object"
}
