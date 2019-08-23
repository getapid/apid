package step

type Validator interface {
	validate(ExpectedResponse, HTTPResponse) ValidationResult
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

func (v *ResponseValidator) validate(ExpectedResponse, HTTPResponse) ValidationResult {
	return ValidationResult{
		OK:     true,
		Errors: make(map[string]string),
	}
}
