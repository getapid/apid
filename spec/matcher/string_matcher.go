package matcher

import (
	"fmt"
	"strings"

	"github.com/getapid/apid/log"
)

type stringMatcher struct {
	value         string `json:"value"`
	caseSensitive bool   `json:"case_sensitive"`
}

func StringMatcher(value string) Matcher {
	return &stringMatcher{
		value:         value,
		caseSensitive: true,
	}
}

func StringMatcherWithOptions(options interface{}) Matcher {
	m, err := extractMap(options)
	if err != nil {
		log.L.Fatalf("invalid string matcher with options %v", options)
		return nil
	}

	matcher := stringMatcher{}

	value, err := extractString(m["value"])
	if err != nil {
		switch err {
		case ErrInterfaceNil:
			log.L.Fatalf("missing string matcher value with options %v", m)
		default:
			log.L.Fatalf("invalid string matcher value with options %v", m)
		}
		return nil
	}
	matcher.value = value

	caseSensitive, err := extractBool(m["case_sensitive"])
	if err != nil {
		switch err {
		case ErrInterfaceNil:
			break
		default:
			log.L.Fatalf("invalid string matcher case_sensitive with options %v", m)
			return nil
		}
	}

	matcher.caseSensitive = caseSensitive
	return matcher
}

func (m stringMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch s := data.(type) {
	case string:
		if m.caseSensitive {
			if strings.Compare(s, m.value) == 0 {
				return true, []string{fmt.Sprintf("%s is %s", location, m.value)}, nil
			}
			return false, nil, []string{fmt.Sprintf("%s: wanted %s, got %v", location, m.value, data)}
		} else {
			if strings.EqualFold(s, m.value) {
				return true, []string{fmt.Sprintf("%s is %s", location, m.value)}, nil
			}
			return false, nil, []string{fmt.Sprintf("%s: wanted %s, got %v", location, m.value, data)}
		}
	default:
		return false, nil, []string{fmt.Sprintf("%s: wanted %s, got %v", location, m.value, data)}
	}
}

func (m stringMatcher) String() string {
	return m.value
}
