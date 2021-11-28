package matcher

import (
	"encoding/json"
	"fmt"

	"github.com/getapid/apid/log"
)

type jsonMatcher struct {
	matchers map[Matcher]Matcher
	subset   bool
}

func mapMatchers(json map[string]interface{}) map[Matcher]Matcher {
	matchers := make(map[Matcher]Matcher, len(json))
	for name, value := range json {
		nameMatcher := isShorthandMatcher(name)
		if nameMatcher == nil {
			nameMatcher = StringMatcher(name)
		}
		matchers[nameMatcher] = GetMatcher(value)
	}
	return matchers
}

func JSONMatcher(json map[string]interface{}) Matcher {
	return &jsonMatcher{
		matchers: mapMatchers(json),
		subset:   true,
	}
}

func JSONMatcherWithOptions(i interface{}) Matcher {
	params, err := extractMap(i)
	if err != nil {
		log.L.Fatal("invalid json matcher params, %v", i)
	}
	subset := optionalBool(params["subset"], false)
	json, err := extractMap(params["value"])
	if err != nil {
		log.L.Fatalf("json matcher missing value, got %v", params)
	}

	return &jsonMatcher{
		subset:   subset,
		matchers: mapMatchers(json),
	}
}

func (m jsonMatcher) Match(data interface{}, location string) (ok bool, pass []string, fail []string) {
	switch d := data.(type) {
	case map[string]interface{}:
		keysMatched := make(map[string]bool, len(d))
		for nameMatcher, valueMatcher := range m.matchers {
			found := false
			errStr := fmt.Sprintf("%s.%s not found", location, nameMatcher)
			for name, value := range d {
				key := fmt.Sprintf("`%s.%s`", location, nameMatcher)
				if ok, _, _ := nameMatcher.Match(name, key); !ok {
					continue
				}

				if ok, _, _ := valueMatcher.Match(value, key); ok {
					pass = append(pass, fmt.Sprintf("%s is %s", key, valueMatcher))
					found = true
					keysMatched[name] = true
					break
				} else {
					errStr = fmt.Sprintf("%s is not %s, got %v", key, valueMatcher, value)
				}
			}
			if !found {
				fail = append(fail, errStr)
			}
		}

		if !m.subset {
			if len(keysMatched) != len(d) {
				var extra []string
				for key := range d {
					if _, ok := keysMatched[key]; !ok {
						extra = append(extra, key)
					}
				}
				fail = append(fail, fmt.Sprintf("%s has extra keys %+q", location, extra))
			}
		}

		return len(fail) == 0, pass, fail
	default:
		return false, nil, []string{fmt.Sprintf("%s expected map, got %v", location, data)}
	}
}

func (m jsonMatcher) String() string {
	result := make(map[string]string, len(m.matchers))
	for name, value := range m.matchers {
		result[name.String()] = value.String()
	}
	val, _ := json.Marshal(result)
	return string(val)
}
