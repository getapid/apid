package matcher

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type intMatcher struct {
	value int64
}

func IntMatcherWithOptions(params interface{}) Matcher {
	switch v := params.(type) {
	case float64:
		// check if value is int, json unmarshal puts
		// numbers in float64
		if v == float64(int(v)) {
			return intMatcher{int64(v)}
		}
		log.L.Fatalf("invalid int matcher, got %v", params)
	default:
		log.L.Fatalf("invalid int matcher, got %v", params)
	}
	return nil
}

func (m intMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch val := data.(type) {
	case int64:
		return m.checkInt(val, location)
	case int:
		return m.checkInt(int64(val), location)
	case float64:
		if val == float64(int(val)) {
			return m.checkInt(int64(val), location)
		}
	}

	return false, nil, []string{fmt.Sprintf("%s: wanted %d, got %s", location, m.value, data)}
}

func (m intMatcher) checkInt(val int64, location string) (bool, []string, []string) {
	if m.value == val {
		return true, []string{fmt.Sprintf("%s is %d", location, m.value)}, nil
	}
	return false, nil, []string{fmt.Sprintf("%s: wanted %d, got %d", location, m.value, val)}
}

func (m intMatcher) String() string {
	return fmt.Sprintf("%d", m.value)
}
