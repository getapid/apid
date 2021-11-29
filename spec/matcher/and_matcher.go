package matcher

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type andMatcher struct {
	matchers []Matcher
}

func AndMatcherWithOptions(params interface{}) Matcher {
	switch v := params.(type) {
	case []interface{}:
		var matchers []Matcher
		for _, value := range v {
			matchers = append(matchers, GetMatcher(value))
		}
		return &andMatcher{matchers: matchers}
	default:
		log.L.Fatalf("invalid and matcher values (expected array), got %v", params)
	}
	return nil
}

func (m andMatcher) Match(data interface{}, location string) (ok bool, pass []string, fail []string) {
	for _, matcher := range m.matchers {
		_, p, f := matcher.Match(data, location)
		pass = append(pass, p...)
		fail = append(fail, f...)
	}
	ok = len(fail) == 0
	return
}

func (m andMatcher) String() string {
	return fmt.Sprintf("%v", m.matchers)
}
