package matcher

import (
	"fmt"
	"regexp"

	"github.com/getapid/apid/log"
)

type regexMatcher struct {
	pattern string
	regexp  *regexp.Regexp
}

func RegexMatcherWithOptions(options interface{}) Matcher {
	p, err := extractString(options)
	if err != nil {
		log.L.Fatalf("invalid regex matcher with options %v", options)
		return nil
	}
	r, err := regexp.Compile(p)
	if err != nil {
		log.L.Fatalf("invalid regex expression %s", p)
		return nil
	}
	return regexMatcher{p, r}
}

func (m regexMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	switch s := data.(type) {
	case string:
		if m.regexp.MatchString(s) {
			return true, []string{fmt.Sprintf("%s matches %s", location, m.pattern)}, nil
		}
		return false, nil, []string{fmt.Sprintf("%s: wanted %s, got %v", location, m.pattern, data)}
	}
	return false, nil, []string{fmt.Sprintf("%s: wanted %s, got %v", location, m.pattern, data)}
}

func (m regexMatcher) String() string {
	return m.pattern
}
