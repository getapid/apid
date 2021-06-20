package step

import (
	"encoding/json"
	"fmt"
	http2 "net/http"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/getapid/cli/common/http"
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

	appendErr(errMsgs, "code", v.validateCode(exp.Code, actual.StatusCode))
	appendErr(errMsgs, "headers", v.validateHeaders(exp.Headers, actual.Header))
	for _, body := range exp.Body {
		appendErr(errMsgs, "body", v.validateBody(body, actual.Body))
	}

	result.Errors = errMsgs
	return
}

func (httpValidator) validateCode(exp *int, actual int) error {
	if exp != nil && *exp != actual {
		return fmt.Errorf("want %d, received %d", *exp, actual)
	}
	return nil
}

func (httpValidator) validateHeaders(exp *Headers, actual http2.Header) error {
	if exp == nil {
		return nil
	}

	headersEqual := func(expected []string, actual []string) bool {
		for _, expVal := range expected {
			found := false
			for _, actualVal := range actual {
				if expVal == actualVal {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}

	var accumulatedErrs error

	for h, expVal := range *exp {
		received, ok := actual[h]
		if !ok {
			err := fmt.Errorf("%q not present in response", h)
			accumulatedErrs = multierr.Append(accumulatedErrs, err)
			continue
		}

		if !headersEqual(expVal, received) {
			err := fmt.Errorf("%q: want %q, received %q", h, received, expVal)
			accumulatedErrs = multierr.Append(accumulatedErrs, err)
		}
	}
	return accumulatedErrs
}

func (httpValidator) validateBody(exp *ExpectBody, actual []byte) error {
	if exp == nil {
		return nil
	}

	var expected interface{}
	err := json.Unmarshal([]byte(exp.Is), &expected)
	if err != nil {
		expected = exp.Is
	}

	message := "expected value for body"
	if exp.Selector != nil {
		message = fmt.Sprintf("expected value for `%s` field", *exp.Selector)
	}

	var received interface{}
	if exp.Selector != nil {
		val := gjson.GetBytes(actual, *exp.Selector)
		if !val.Exists() {
			return fmt.Errorf("invalid selector: %s", *exp.Selector)
		}
		received = val.Value()
	} else {
		err = json.Unmarshal(actual, &received)
		if err != nil {
			if !plainTextEqual(exp.Is, string(actual), *exp.Subset) {
				return fmt.Errorf("%s doesn't match actual:\nwant:\n%s\nreceived:\n%s", message, prettyPrint(exp.Is), prettyPrint(actual))
			}
			return nil
		}
	}

	if !mapStructsEqual(expected, received, *exp.KeysOnly, *exp.Subset) {
		return fmt.Errorf("%s doesn't match actual:\nwant:\n%s\nreceived:\n%s", message, prettyPrint(expected), prettyPrint(received))
	}
	return nil
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func plainTextEqual(exp, actual string, subset bool) bool {
	if !subset {
		return exp == actual
	}
	return strings.Index(actual, exp) != -1
}

// mapStructsEqual checks if the actual map has all the fields and their corresponding values that exp has.
// If the value for a field in exp is also a map, then mapStructsEqual will check recursively there as well
// (returning false if the type in actual is not another map)
func mapStructsEqual(exp, actual interface{}, keysOnly, subset bool) bool {
	switch exp.(type) {
	case map[string]interface{}:
		expMap := exp.(map[string]interface{})
		actualMap, ok := actual.(map[string]interface{})
		if !ok {
			return false
		}
		if !subset {
			// check if all the keys in the actual map are in the expected map
			for k := range actualMap {
				if _, ok := expMap[k]; !ok {
					return false
				}
			}
		}
		for k, expNested := range expMap {
			if actualNested, ok := actualMap[k]; !ok {
				return false
			} else {
				if !mapStructsEqual(expNested, actualNested, keysOnly, subset) {
					return false
				}
			}
		}
	case []interface{}:
		expSlice := exp.([]interface{})
		actualSlice, ok := actual.([]interface{})
		if !ok {
			return false
		}

		if !subset && len(expSlice) != len(actualSlice) {
			return false
		}

		if len(expSlice) == 0 || !ok {
			return ok
		}

		for _, expVal := range expSlice {
			found := false
			for _, actualVal := range actualSlice {
				if mapStructsEqual(expVal, actualVal, keysOnly, subset) {
					found = true
				}
			}
			if !found {
				return false
			}
		}
	default:
		if !keysOnly {
			return exp == actual
		}
	}
	return true
}
