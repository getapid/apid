package step

import "github.com/iv-p/apid/common/http"

type validator interface {
	validate(ExpectedResponse, *http.Response) (ValidationResult, error)
}

type httpValidator struct{}

// ValidationResult holds information if the validation succeeded or not and what
// errors were encountered if any
type ValidationResult struct {
	OK     bool              // overall check status, true only if every other check passes
	Errors map[string]string // a list of error keys and more information about what caused them
}

// NewHTTPValidator instantiates a new HTTPValidator
func NewHTTPValidator() validator {
	return &httpValidator{}
}

func (v *httpValidator) validate(ExpectedResponse, *http.Response) (ValidationResult, error) {
	return ValidationResult{
		OK:     true,
		Errors: make(map[string]string),
	}, nil
}
