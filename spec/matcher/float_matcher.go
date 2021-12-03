package matcher

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type floatMatcher struct {
	value float64
}

func FloatMatcherWithOptions(value interface{}) Matcher {
	switch v := value.(type) {
	case float64:
		return floatMatcher{v}
	default:
		log.L.Fatalf("invalid float matcher, got %v", value)
	}
	return nil
}

func (m floatMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch val := data.(type) {
	case float64:
		if m.value == val {
			return true, []string{fmt.Sprintf("%s is %.f", location, m.value)}, nil
		}
	case int64:
		if m.value == float64(val) {
			return true, []string{fmt.Sprintf("%s is %f", location, m.value)}, nil
		}
	case int:
		if m.value == float64(val) {
			return true, []string{fmt.Sprintf("%s is %f", location, m.value)}, nil
		}
	}

	return false, nil, []string{fmt.Sprintf("%s: wanted %f, got %v", location, m.value, data)}
}

func (m floatMatcher) String() string {
	return fmt.Sprintf("%f", m.value)
}
