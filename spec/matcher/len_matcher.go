package matcher

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type lenMatcher struct {
	value int
}

func LenMatcherWithOptions(params interface{}) Matcher {
	switch v := params.(type) {
	case float64:
		// check if value is int, json unmarshal puts
		// numbers in float64
		if v == float64(int(v)) {
			return lenMatcher{int(v)}
		}
		log.L.Fatalf("invalid len matcher, got %v", params)
	default:
		log.L.Fatalf("invalid len matcher, got %v", params)
	}
	return nil
}

func (m lenMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch val := data.(type) {
	case []interface{}:
		if len(val) == m.value {
			return true, []string{fmt.Sprintf("`%s` has length %d", location, m.value)}, nil
		}
		return true, nil, []string{fmt.Sprintf("`%s` of length %d wanted, got %d", location, m.value, len(val))}
	case map[string]interface{}:
		if len(val) == m.value {
			return true, []string{fmt.Sprintf("`%s` has length %d", location, m.value)}, nil
		}
		return true, nil, []string{fmt.Sprintf("`%s` of length %d wanted, got %d", location, m.value, len(val))}
	case string:
		if len(val) == m.value {
			return true, []string{fmt.Sprintf("`%s` has length %d", location, m.value)}, nil
		}
		return false, nil, []string{fmt.Sprintf("`%s` of length %d wanted, got %d", location, m.value, len(val))}
	default:
		return false, nil, []string{fmt.Sprintf("`%s` not a map, array or string", location)}
	}

}

func (m lenMatcher) String() string {
	return fmt.Sprintf("of lentgh %d", m.value)
}
