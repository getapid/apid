package matcher

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type boolMatcher struct {
	value bool
}

func BoolMatcherWithOptions(params interface{}) Matcher {
	switch v := params.(type) {
	case bool:
		return boolMatcher{v}
	default:
		log.L.Fatalf("invalid bool matcher, got %v", params)
	}
	return nil
}

func (m boolMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch val := data.(type) {
	case bool:
		if m.value == val {
			return true, []string{fmt.Sprintf("%s is %t", location, m.value)}, nil
		}
	}

	return false, nil, []string{fmt.Sprintf("%s: wanted %t, got %v", location, m.value, data)}
}

func (m boolMatcher) String() string {
	return fmt.Sprintf("%t", m.value)
}
