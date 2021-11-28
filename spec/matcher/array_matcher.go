package matcher

import (
	"encoding/json"
	"fmt"

	"github.com/getapid/apid/log"
)

type arrayMatcher struct {
	matchers []Matcher
	subset   bool
}

func sliceMatchers(slice []interface{}) []Matcher {
	var matchers []Matcher
	for _, value := range slice {
		matchers = append(matchers, GetMatcher(value))
	}
	return matchers
}

func ArrayMatcher(i []interface{}) Matcher {
	return &arrayMatcher{
		matchers: sliceMatchers(i),
		subset:   true,
	}
}

func ArrayMatcherWithOptions(i interface{}) Matcher {
	params, err := extractMap(i)
	if err != nil {
		log.L.Fatal("invalid array matcher params, %v", i)
	}
	subset := optionalBool(params["subset"], false)
	slice, err := extractSlice(params["value"])
	if err != nil {
		log.L.Fatalf("array matcher missing value, got %v", params)
	}

	return &arrayMatcher{
		subset:   subset,
		matchers: sliceMatchers(slice),
	}
}

func (m arrayMatcher) Match(data interface{}, location string) (ok bool, pass []string, fail []string) {
	key := fmt.Sprintf("`%s`", location)
	switch d := data.(type) {
	case []interface{}:
		keysMatched := make(map[int]bool, len(d))
		for i, matcher := range m.matchers {
			found := false
			for _, value := range d {
				ok, _, _ := matcher.Match(value, location)
				if ok {
					pass = append(pass, fmt.Sprintf("%s is %s", key, matcher))
					found = true
					keysMatched[i] = true
				}
			}
			if !found {
				fail = append(fail, fmt.Sprintf("%s elements do not match %s", key, matcher))
			}
		}
		if !m.subset {
			if len(d) > len(keysMatched) {
				fail = append(fail, fmt.Sprintf("%s has extra keys", key))
			}
		}
		return len(fail) == 0, pass, fail
	default:
		return false, nil, []string{fmt.Sprintf("%s expected array, got %v", key, data)}
	}
}

func (m arrayMatcher) String() string {
	var result []string
	for _, value := range m.matchers {
		result = append(result, value.String())
	}
	val, _ := json.Marshal(result)
	return string(val)
}
