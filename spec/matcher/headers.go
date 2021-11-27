package matcher

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type HeaderMatcher map[StringMatcher]StringMatcher

func (v *HeaderMatcher) Set(headers map[string]string) {
	m := make(map[StringMatcher]StringMatcher, len(headers))
	for name, value := range headers {
		m[StringMatcher(name)] = StringMatcher(value)
	}
	val := HeaderMatcher(m)
	v = &val
}

func (v HeaderMatcher) Validate(headers map[string][]string) (pass []string, fail []string) {
	log.L.Infof("validating headers %v", headers)
	for nameMatcher, valueMatcher := range v {
		found := false
		errStr := fmt.Sprintf("expected header %s, but none found", nameMatcher)
		for name, values := range headers {
			if !nameMatcher.Validate(name) {
				continue
			}
			errStr = fmt.Sprintf("expected header %s to match %s, got %v", nameMatcher, valueMatcher, values)

			for _, value := range values {
				if !valueMatcher.Validate(value) {
					continue
				}
				pass = append(pass, fmt.Sprintf("has header %s = %s", nameMatcher, valueMatcher))
				found = true
				break
			}
			if found {
				break
			}
		}
		if !found {
			fail = append(fail, errStr)
		}
	}
	return pass, fail
}
