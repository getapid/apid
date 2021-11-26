package step

import (
	"github.com/getapid/apid/http"
	"github.com/getapid/apid/spec"
)

type Validator struct{}

// NewValidator instantiates a new Validator
func NewValidator() Validator {
	return Validator{}
}

func (v Validator) validate(exp spec.Expect, actual *http.Response) (ok bool, passed []string, failed []string) {
	if exp.Code != nil {
		pass, fail := exp.Code.Validate(int64(actual.StatusCode))
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

	if exp.Text != nil {
		pass, fail := exp.Text.Validate(string(actual.Body))
		failed = append(failed, fail...)
		passed = append(passed, pass...)
	}

	if len(exp.JSON) > 0 {
		for _, validator := range exp.JSON {
			pass, fail := validator.Validate(actual.Body)
			failed = append(failed, fail...)
			passed = append(passed, pass...)
		}
	}

	ok = len(failed) == 0
	return
}
