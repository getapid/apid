package matcher

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type rangeMatcher struct {
	from, to float64
}

func RangeMatcherWithOptions(value interface{}) Matcher {
	params, err := extractMap(value)
	if err != nil {
		log.L.Fatalf("invalid range matcher, got %v", value)
	}

	from, err := extractFloat(params["from"])
	if err != nil {
		log.L.Fatalf("invalid range matcher, got %v", value)
	}

	to, err := extractFloat(params["to"])
	if err != nil {
		log.L.Fatalf("invalid range matcher, got %v", value)
	}

	return &rangeMatcher{from: from, to: to}
}

func (m rangeMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch val := data.(type) {
	case float64:
		return m.check(val, location)
	case int64:
		return m.check(float64(val), location)
	}

	return false, nil, []string{fmt.Sprintf("%s wanted in range [%f, %f], got %v", location, m.from, m.to, data)}
}

func (m rangeMatcher) check(val float64, location string) (bool, []string, []string) {
	if val >= m.from && val <= m.to {
		return true, []string{fmt.Sprintf("%s is between %f and %f", location, m.from, m.to)}, nil
	}
	return false, nil, []string{fmt.Sprintf("%s is not between %f and %f", location, m.from, m.to)}
}

func (m rangeMatcher) String() string {
	return fmt.Sprintf("[%f, %f]", m.from, m.to)
}
