package step

import (
	"fmt"
	http2 "net/http"
	"strings"

	"github.com/iv-p/apid/common/http"
	"go.uber.org/multierr"
)

type validator interface {
	validate(ExpectedResponse, *http.Response) ValidationResult
}

type httpValidator struct{}

// ValidationResult holds information if the validation succeeded or not and what
// errors were encountered if any
type ValidationResult struct {
	Errors map[string]string // a list of error keys and more information about what caused them
}

// OK returns overall check status, true only if every other check passes
func (r ValidationResult) OK() bool {
	return len(r.Errors) == 0
}

// NewHTTPValidator instantiates a new HTTPValidator
func NewHTTPValidator() validator {
	return httpValidator{}
}

func (v httpValidator) validate(exp ExpectedResponse, actual *http.Response) (result ValidationResult) {
	errMsgs := make(map[string]string)
	appendErr := func(errors map[string]string, key string, err error) {
		if err != nil {
			errors[key] = err.Error()
		}
	}

	appendErr(errMsgs, "code", v.checkCode(exp.Code, actual.StatusCode))
	appendErr(errMsgs, "headers", v.checkHeaders(exp.Headers, actual.Header))

	result.Errors = errMsgs
	return
}

func (httpValidator) checkCode(exp *int, actual int) error {
	if exp != nil && *exp != actual {
		return fmt.Errorf("want %d, received %d", *exp, actual)
	}
	return nil
}

func (httpValidator) checkHeaders(exp *Headers, actual http2.Header) error {
	if exp == nil {
		return nil
	}

	var accumulatedErrs error

	for h, expVal := range *exp {
		received, ok := actual[h]
		if !ok {
			err := fmt.Errorf("%q not present in response", h)
			accumulatedErrs = multierr.Append(accumulatedErrs, err)
			continue
		}

		actualVal := strings.Join(received, "")
		if expVal != actualVal {
			err := fmt.Errorf("%q: want %q, received %q", h, actualVal, expVal)
			accumulatedErrs = multierr.Append(accumulatedErrs, err)
		}
	}
	return accumulatedErrs
}
