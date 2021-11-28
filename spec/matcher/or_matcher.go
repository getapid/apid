package matcher

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type orMatcher struct {
	matchers []Matcher
}

func OrMatcherWithOptions(params interface{}) Matcher {
	switch v := params.(type) {
	case []interface{}:
		var matchers []Matcher
		for _, value := range v {
			matchers = append(matchers, GetMatcher(value))
		}
		return &orMatcher{matchers: matchers}
	default:
		log.L.Fatalf("invalid or matcher values (expected array), got %v", params)
	}
	return nil
}

func (m orMatcher) Match(data interface{}, location string) (ok bool, pass []string, fail []string) {
	for _, matcher := range m.matchers {
		if pa, p, _ := matcher.Match(data, location); pa {
			pass = append(pass, p...)
			ok = true
			return
		}
	}
	ok = false
	fail = append(fail, fmt.Sprintf("`%s` did not meet any criteria %s", location, m.matchers))
	return
}

func (m orMatcher) String() string {
	return fmt.Sprintf("%d", m.matchers)
}
