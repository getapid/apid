package step

import "github.com/iv-p/apid/common/http"

type Validator interface {
	validate(ExpectedResponse, *http.Response) ValidationResult
}

type ResponseValidator struct {
	Validator
}

type ValidationResult struct {
	OK     bool              // overall check status, true only if every other check passes
	Errors map[string]string // a list of error keys and more information about what caused them
}

func NewResponseValidator() Validator {
	return &ResponseValidator{}
}

func (v *ResponseValidator) validate(ExpectedResponse, *http.Response) ValidationResult {
	return ValidationResult{
		OK:     true,
		Errors: make(map[string]string),
	}
}
