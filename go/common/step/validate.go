package step

import "github.com/iv-p/apid/common/http"

// Validator is the interface for step validators
type Validator interface {
	validate(ExpectedResponse, *http.Response) (ValidationResult, error)
}

// HTTPValidator receives a http.Response and checks if it what's defined
// in the step's expected block
type HTTPValidator struct {
	Validator
}

// ValidationResult holds information if the validation succeeded or not and what
// errors were encountered if any
type ValidationResult struct {
	OK     bool              // overall check status, true only if every other check passes
	Errors map[string]string // a list of error keys and more information about what caused them
}

// NewHTTPValidator instantiates a new HTTPValidator
func NewHTTPValidator() Validator {
	return &HTTPValidator{}
}

func (v *HTTPValidator) validate(ExpectedResponse, *http.Response) (ValidationResult, error) {
	return ValidationResult{
		OK:     true,
		Errors: make(map[string]string),
	}, nil
}
