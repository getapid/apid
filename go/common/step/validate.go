package step

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	http2 "net/http"
	"reflect"
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

func (httpValidator) validateBody(exp *ExpectBody, actual io.Reader) error {
	const (
		typeJson  = "json"
		typePlain = "plaintext"
	)

	if exp == nil {
		return nil
	}

	typ := typePlain
	if exp.Type != nil {
		typ = *exp.Type
	}
	exact := true
	if exp.Exact != nil {
		exact = *exp.Exact
	}

	if (typ != typeJson && typ != typePlain) && exact {
		return fmt.Errorf(`cannot check non-exact body with type %q, only "json" supported`, typ)
	}

	var unmarshall func([]byte, interface{}) error
	//typ := "plaintext" // todo use this

	switch typ {
	case typeJson:
		unmarshall = json.Unmarshal
	case typePlain:
		unmarshall = func([]byte, interface{}) error {
			panic("implement me")
		}
	default:
		return fmt.Errorf("no support for type %q", *exp.Type)
	}

	var expected interface{}
	err := unmarshall([]byte(exp.Content), &expected) // todo
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
		return fmt.Errorf("coulnd't convert response to type %q, response: %s", typ, body) // TODO remove this dereference here and use the type
	}
	// todo interpolate the expected body

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
	return true // TODO
}
