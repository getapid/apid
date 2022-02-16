package matcher

import (
	"fmt"
)

type typeBoolMatcher struct{}

func TypeBoolMatcherWithOptions(value interface{}) Matcher {
	return &typeBoolMatcher{}
}

func (m typeBoolMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch data.(type) {
	case bool:
		return true, []string{fmt.Sprintf("%s is bool", location)}, nil
	}

	return false, nil, []string{fmt.Sprintf("%s wanted bool, got %v", location, data)}
}

func (m typeBoolMatcher) String() string {
	return "type::bool"
}
