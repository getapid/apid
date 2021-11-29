package matcher

import (
	"fmt"
)

type typeStringMatcher struct{}

func TypeStringMatcherWithOptions(value interface{}) Matcher {
	return &typeStringMatcher{}
}

func (m typeStringMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch data.(type) {
	case string:
		return true, []string{fmt.Sprintf("%s is string", location)}, nil
	}

	return false, nil, []string{fmt.Sprintf("%s wanted string, got %v", location, data)}
}

func (m typeStringMatcher) String() string {
	return "type::string"
}
