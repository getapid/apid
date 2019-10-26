package step

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	http2 "net/http"
	"reflect"

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
	defer actual.Body.Close()
	errMsgs := make(map[string]string)
	appendErr := func(errors map[string]string, key string, err error) {
		if err != nil {
			errors[key] = err.Error()
		}
	}

	appendErr(errMsgs, "code", v.validateCode(exp.Code, actual.StatusCode))
	appendErr(errMsgs, "headers", v.validateHeaders(exp.Headers, actual.Header))
	appendErr(errMsgs, "body", v.validateBody(exp.Body, actual.Body))

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

func (httpValidator) validateBody(exp *ExpectBody, actual io.Reader) error {
	if exp == nil {
		return nil
	}

	// validation should have set those to non-nil
	typ := *exp.Type
	exact := *exp.Exact

	var unmarshall func([]byte, interface{}) error

	switch typ {
	case "json":
		unmarshall = json.Unmarshal
	case "plaintext":
		unmarshall = func(b []byte, v interface{}) error {
			reflect.ValueOf(v).Elem().Set(reflect.ValueOf(string(b)))
			return nil
		}
	default: // again, should have been covered by validation
		panic(fmt.Errorf("no support for type %q", typ))
	}

	var expected interface{}
	err := unmarshall([]byte(exp.Content), &expected)
	if err != nil {
		return fmt.Errorf("couldn't convert expected body into type: %w, body = %s", err, exp.Content)
	}

	var received interface{}
	body, err := ioutil.ReadAll(actual)
	if err != nil {
		return err
	}
	err = unmarshall(body, &received)
	if err != nil {
		return fmt.Errorf("coulnd't convert response to type %q, response: %s", typ, body)
	}

	if exact {
		if !reflect.DeepEqual(expected, received) {
			return fmt.Errorf("expected body doesn't match actual: want = %#v, received = %#v", expected, received)
		}
	} else {
		if !fieldsEqual(expected, received) {
			return fmt.Errorf("expected body's fields don't match actual: want = %#v, received = %#v", expected, received)
		}
	}
	return nil
}

func fieldsEqual(exp, actual interface{}) bool {
	switch expMap := exp.(type) {
	case map[string]interface{}:
		actualMap, ok := actual.(map[string]interface{})
		if !ok {
			return false
		}
		for k, expNested := range expMap {
			if actualNested, ok := actualMap[k]; !ok {
				return false
			} else {
				if !fieldsEqual(expNested, actualNested) {
					return false
				}
			}
		}
	}
	return true
}
