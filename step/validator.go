package step

import (
	"github.com/getapid/apid/http"
	"github.com/getapid/apid/spec"
)

type Matcher struct{}

// NewMatcher instantiates a new Matcher
func NewMatcher() Matcher {
	return Matcher{}
}

func (v Matcher) validate(exp spec.Expect, actual *http.Response) (ok bool, passed []string, failed []string) {
	if exp.Code != nil {
		pass, fail := exp.Code.Validate(actual.StatusCode)
		failed = append(failed, fail...)
		passed = append(passed, pass...)
	}

	if exp.Headers != nil {
		actualHeaders := make(map[string][]string, len(actual.Header))
		for name, values := range actual.Header {
			actualHeaders[name] = values
		}
		pass, fail := exp.Headers.Validate(actualHeaders)
		failed = append(failed, fail...)
		passed = append(passed, pass...)
	}

	if exp.Body != nil {
		pass, fail := exp.Body.Validate(actual.Body)
		failed = append(failed, fail...)
		passed = append(passed, pass...)
	}

	// if len(exp.JSON) > 0 {
	// 	for _, matcher := range exp.JSON {
	// 		pass, fail := matcher.Validate(actual.Body)
	// 		failed = append(failed, fail...)
	// 		passed = append(passed, pass...)
	// 	}
	// }

	ok = len(failed) == 0
	return
}
